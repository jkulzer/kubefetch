package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	//command-line flags
	"flag"

	//allows converting from ints to strings
	"strconv"

	//allows the ascii art to be stored in the binary itself
	"embed"
	"io/fs"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// this embeds the assets in the 'assets' folders
//
//go:embed assets
var asciiArtFile embed.FS

func main() {

	versionFlag := flag.Bool("version", false, "Print the version")
	flag.Parse()

	if *versionFlag {
		fmt.Println("0.5.3")
		return
	} else {

		printArt()

	}
}

func getKubeconfig() (*rest.Config, error) {
	// Get the path to the kubeconfig file
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	// Build the client config from the kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func getPodCount(podCount *int) {

	// create the clientset
	config, err := getKubeconfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//for {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// Get the number of pods
	*podCount = len(pods.Items)
}

func getNamespaceCount(namespaceCount *int) {

	// create the clientset
	config, err := getKubeconfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// Get the number of pods
	*namespaceCount = len(namespaces.Items)
}

func getNodeCount() (int64, int64) {

	// create the clientset
	config, err := getKubeconfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// Get the number of nodes
	nodeCount := len(nodes.Items)

	// Get the total capacity of each node in terms of pods
	podsPerNode := make([]int, nodeCount)
	for i, node := range nodes.Items {
		podsPerNode[i] = int(node.Status.Capacity.Pods().Value())
	}

	// Get the maximum number of pods in the cluster
	// this is in the nodeCount function and not in the podCount function because it depends on the count of the nodes
	// it gets the amount of nodes and loops through every node to get the amount of pods available on the node
	var maxPods int
	for _, pods := range podsPerNode {
		maxPods += pods
	}

	return int64(nodeCount), int64(maxPods)
}

func getKubeVersion() (string, error) {

	// create the clientset
	config, err := getKubeconfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Get the server version
	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}

	return version.String(), nil
}

func getContainerRuntimeInterface() (string, error) {

	// create the clientset
	config, err := getKubeconfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Retrieve the CRI information from the node status
	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	if len(nodeList.Items) == 0 {
		return "", fmt.Errorf("no nodes found")
	}

	cri := nodeList.Items[0].Status.NodeInfo.ContainerRuntimeVersion

	return cri, nil

}

func getStorage() (string, error) {

	// create the clientset
	config, err := getKubeconfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Retrieve the CRI information from the node status
	storageClassList, err := clientset.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	// Check if "longhorn-system" namespace exists
	storageUsed := "No Storage detected"
	for _, storageclass := range storageClassList.Items {
		if storageclass.Name == "longhorn" {
			storageUsed = "Longhorn"
		} else if strings.Contains(storageclass.Name, "rook") {
			storageUsed = "Rook/Ceph"
		} else {
			storageUsed = storageUsed + ""
		}
	}

	return storageUsed, nil

}

func getKubernetesEndpointPort() int {
	// create the clientset
	config, err := getKubeconfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	endpoint, err := clientset.CoreV1().Endpoints("default").Get(context.TODO(), "kubernetes", v1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	var portNumber int
	// Accessing the ports exposed by the Endpoint
	for _, subset := range endpoint.Subsets {
		for _, port := range subset.Ports {
			// Accessing port.Port
			portNumber = int(port.Port)
		}
	}

	return portNumber
}

func getGitops() (string, error) {

	// create the clientset
	config, err := getKubeconfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Retrieve the list of namespaces to find out if Flux or ArgoCD namespaces exist
	namespaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	var hasFlux, hasArgoCD bool // Flags to track namespace detection

	for _, namespace := range namespaceList.Items {
		if namespace.Name == "flux-system" {
			hasFlux = true
		} else if strings.Contains(namespace.Name, "argocd") {
			hasArgoCD = true
		}
	}

	var gitopsToolUsed string

	if hasFlux && hasArgoCD {
		gitopsToolUsed = "Argo CD + Flux"
	} else if hasFlux {
		gitopsToolUsed = "Flux"
	} else if hasArgoCD {
		gitopsToolUsed = "Argo CD"
	} else {
		gitopsToolUsed = "No GitOps tool used"
	}

	return gitopsToolUsed, nil
}

func getNodeAge() (int64, error) {

	// create the clientset
	config, err := getKubeconfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	var oldestNode *corev1.Node
	var oldestNodeAge time.Time

	for _, node := range nodes.Items {
		creationTimestamp := node.CreationTimestamp.Time
		if oldestNode == nil || creationTimestamp.Before(oldestNodeAge) {
			oldestNode = &node
			oldestNodeAge = creationTimestamp
		}
	}
	var clusterAge time.Duration

	clusterAge = time.Now().Sub(oldestNodeAge)
	clusterAgeInDays := int64(clusterAge.Hours() / 24)

	return clusterAgeInDays, nil
}

func printArt() {
	//gets kubernetes version
	version, err := getKubeVersion()
	if err != nil {
		panic(err.Error())
	}

	kubernetesEndpointPort := getKubernetesEndpointPort()

	//gets kubernetes distro
	var distro string
	if strings.Contains(version, "k3s") {
		distro = "K3s"
	} else {

		//microk8s usually uses the apiServer port 16443
		if kubernetesEndpointPort == 16443 {
			distro = "MicroK8s"
		} else {
			distro = "K8s"
		}

	}

	//gets storage solution used
	storage, err := getStorage()
	if err != nil {
		panic(err.Error())
	}

	//gets container runtime interface
	containerRuntimeInterface, err := getContainerRuntimeInterface()
	if err != nil {
		panic(err.Error())
	}

	gitopsTool, err := getGitops()
	if err != nil {
		panic(err.Error())
	}
	clusterAge, err := getNodeAge()
	if err != nil {
		panic(err.Error())
	}

	nodeCount, maxPods := getNodeCount()
	if err != nil {
		panic(err.Error())
	}

	var asciiArtColor string
	var podCount int
	var namespaceCount int

	getPodCount(&podCount)
	getNamespaceCount(&namespaceCount)
	//fetch all values from the various functions above

	var red int
	var green int
	var blue int
	if distro == "MicroK8s" {
		red = 233
		green = 84
		blue = 32
	} else if distro == "K3s" {
		red = 255
		green = 198
		blue = 28
	} else {
		red = 50
		green = 108
		blue = 229
	}

	// Generate True Color escape sequence
	colorCode := fmt.Sprintf("\x1b[38;2;%d;%d;%dm", red, green, blue)
	resetCode := "\x1b[0m" // Reset color escape sequence

	// Read ASCII art from file
	content, err := fs.ReadFile(asciiArtFile, "assets/"+distro)
	if err != nil {
		fmt.Println("Error reading ascii art, should NEVER happen", err)
		return
	}

	// Split the ASCII art into lines
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")

	// Additional data for each line
	additionalData := []string{
		colorCode + asciiArtColor + "    " + "Distro: " + resetCode + distro,
		colorCode + asciiArtColor + "    " + "Version: " + resetCode + version,
		colorCode + asciiArtColor + "    " + "Node Count: " + resetCode + fmt.Sprint(nodeCount),
		colorCode + asciiArtColor + "    " + "Pod Count: " + resetCode + strconv.Itoa(podCount) + "/" + fmt.Sprint(maxPods),
		colorCode + asciiArtColor + "    " + "Namespace Count: " + resetCode + strconv.Itoa(namespaceCount),
		colorCode + asciiArtColor + "    " + "Container Runtime Interface: " + resetCode + containerRuntimeInterface,
		colorCode + asciiArtColor + "    " + "Storage: " + resetCode + storage,
		colorCode + asciiArtColor + "    " + "GitOps Tool: " + resetCode + gitopsTool,
		colorCode + asciiArtColor + "    " + "Cluster Age: " + resetCode + fmt.Sprint(clusterAge) + "d",
	}

	// Print each line of the ASCII art along with different additional data
	for i, line := range lines {
		if i < len(additionalData) {
			fmt.Println(line, additionalData[i])
		} else {
			fmt.Println(line) // If additional data runs out, print only the line
		}
	}

}
