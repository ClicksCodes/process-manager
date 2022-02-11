package containerManager

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/go-cni"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Options struct {
	id   string
	name string
}

func getAllContainers() {

}

func getContainers(ids []string) {

}

func createContainer(options Options) {

}

func getLatestConfig(id string) string {
	var files []string

	// Get all the files in the config directory
	root := "./containerManager/config"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".nix" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Sort the files by name
	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})

	// Get the latest file
	latest := files[len(files)-1]

	// Return the latest file
	// We return without the extension, as we will add it later
	return strings.TrimSuffix(filepath.Base(latest), filepath.Ext(latest))
}

func GetVersion(id string, version string) string {
	if version == "@latest" {
		log.Println("Finding the latest config")
		version = getLatestConfig(id)
	}
	return version
}

func BuildContainer(id string, version string) error {
	version = GetVersion(id, version)

	// Build with nixos
	// See https://nixos.org/guides/building-and-running-docker-images.html
	log.Println("Building container " + id + " with version " + version)

	log.Println("./containerManager/config/" + id + "/" + version + ".nix")

	cmd := exec.Command("sudo", "nix-build", "./containerManager/config/"+id+"/"+version+".nix")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		log.Println("Error building container: " + err.Error())
		return err
	}
	symlinks, err := filepath.EvalSymlinks("./result")
	if err != nil {
		return err
	}
	// Move the result to the container directory
	log.Println("Moving " + filepath.Base(symlinks) + " to ./containerManager/containers/" + id)
	err = os.Rename("./result", "./containerManager/containers/"+id)
	if err != nil {
		return err
	}
	return nil
}

func RunContainer(id string, version string) error {
	log.Println("Ready to run container " + id + " with version " + version)

	network, err := cni.New()
	if err != nil {
		return err
	}

	if err := network.Load(cni.WithLoNetwork); err != nil {
		return err
	}
	if err := network.Load(cni.WithConfFile("./containerManager/networking/bridge.json")); err != nil {
		return err
	}

	log.Println("Created container network & loaded configuration")

	version = GetVersion(id, version)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	log.Println("Current time: " + timestamp)

	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}

	//defer func(client *containerd.Client) {
	//	err := client.Close(); if err != nil { panic(err) }
	//}(client)

	ctx := namespaces.WithNamespace(context.Background(), "clicks-container-manager")

	log.Println("Created containerd client")

	//IP := net.Interfaces
	// IP := net.Interfaces["eth"].IPConfigs[0].IP.String()
	// fmt.Printf("IP of the default interface %s:%s", "eth0", IP)

	symlinks, err := filepath.EvalSymlinks("./containerManager/containers/" + id)
	if err != nil {
		return err
	}

	log.Println("The selected container has its image at " + symlinks)

	file, err := os.Open(symlinks)
	if err != nil {
		return err
	}

	//defer func(file *os.File) {
	//	err := file.Close(); if err != nil { panic(err) }
	//}(file)

	log.Println("Opened container file for reading")

	reader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}

	log.Println("Ungzipped container file")

	imported, err := client.Import(ctx, reader)
	if err != nil {
		return err
	}
	// See https://blog.scottlowe.org/2020/01/25/manually-loading-container-images-with-containerd/
	image := containerd.NewImage(client, imported[0])

	log.Printf("Successfully imported %s image\n", image.Name())

	if err := image.Unpack(ctx, containerd.DefaultSnapshotter); err != nil {
		return err
	}

	snapshot := containerd.WithNewSnapshot("clicks-container-manager-snapshot-"+id+"-"+version+"-"+timestamp, image)

	container, err := client.NewContainer(
		ctx,
		"clicks-container-manager-"+id+"-"+version+"-"+timestamp,
		containerd.WithImage(image),
		snapshot,
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		return err
	}

	//defer func(container containerd.Container, ctx context.Context, opts ...containerd.DeleteOpts) {
	//	if err := container.Delete(ctx, opts...); err != nil { panic(err) }
	//}(container, ctx)

	net, err := network.Setup(ctx, id+"-"+timestamp, fmt.Sprintf("/proc/%d/ns/net", os.Getpid()))
	if err != nil {
		return err
	}
	// Print out all the interfaces along with their IP addresses
	for key, _ := range net.Interfaces {
		log.Println(key)
	}

	log.Printf("Successfully loaded %s container\n", container.ID())

	if err = mount.SetTempMountLocation("./containerManager/mounts/" + id); err != nil {
		return err
	}

	log.Println("Set mount point to ./containerManager/mounts/" + id)

	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return err
	}

	//defer func(task containerd.Task, ctx context.Context, opts ...containerd.ProcessDeleteOpts) {
	//	if _, err := task.Delete(ctx, opts...); err != nil { panic(err) }
	//}(task, ctx)

	log.Println("Created run-task")
	log.Println(task.Metrics(ctx))

	// Run the container!
	if err := task.Start(ctx); err != nil {
		return err
	}

	status, err := task.Wait(ctx)
	if err != nil {
		return err
	}

	code := <-status

	exitCode, timeToRun, err := code.Result()
	if err != nil {
		return err
	}

	log.Println("Container exited with code " + strconv.Itoa(int(exitCode)) + " at " + timeToRun.String())

	log.Println("Started container!")

	return nil
}

/// See also https://github.com/kubernetes/kubernetes/issues/54918 and add to notes
