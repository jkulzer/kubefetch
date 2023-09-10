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
		fmt.Println("0.7.1")
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
		panic("Failed to fetch kubeconfig. Check if your kubeconfig exists at " + kubeconfig)
	}

	return config, nil
}

func getNonRunningPodCount(clientset *kubernetes.Clientset, pods *corev1.PodList) int {

	var getNonRunningPodCount int
	for _, pod := range pods.Items {

		if pod.Status.Phase == "Running" || pod.Status.Phase == "Succeeded" {
		} else {
			getNonRunningPodCount = getNonRunningPodCount + 1
		}
	}

	return getNonRunningPodCount
}

func getPodCount(clientset *kubernetes.Clientset, pods *corev1.PodList) int {

	// Get the number of pods
	podCount := len(pods.Items)

	return podCount
}

func getServiceCount(clientset *kubernetes.Clientset) int {

	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	serviceCount := len(services.Items)

	return serviceCount

}

func getNamespaceCount(clientset *kubernetes.Clientset) int {

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// Get the number of namespaces
	namespaceCount := len(namespaces.Items)

	return namespaceCount
}

func getNodeCount(clientset *kubernetes.Clientset) (int64, int64) {

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

func getKubeVersion(clientset *kubernetes.Clientset) (string, error) {

	// Get the server version
	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}

	return version.String(), nil
}

func getContainerRuntimeInterface(clientset *kubernetes.Clientset) (string, error) {

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

func getStorage(clientset *kubernetes.Clientset) (string, error) {

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
		} else if strings.Contains(storageclass.Name, "ceph") {
			storageUsed = "Rook/Ceph"
		} else {
			storageUsed = storageUsed + ""
		}
	}

	return storageUsed, nil

}

func getKubernetesEndpointPort(clientset *kubernetes.Clientset) int {

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

func getGitops(clientset *kubernetes.Clientset) (string, error) {

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

func getNodeInfo(clientset *kubernetes.Clientset) *corev1.NodeList {

	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	return nodes
}

func isTalos(nodes *corev1.NodeList) bool {

	var isTalos bool
	for _, node := range nodes.Items {
		nodeInfo := node.Status.NodeInfo.OSImage
		if strings.Contains(nodeInfo, "Talos") {
			isTalos = true
			break
		} else {
			isTalos = false
		}
	}

	return isTalos

}

func getNodeAge(nodes *corev1.NodeList) (int64, error) {

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

func getCNI(clientset *kubernetes.Clientset) string {

	// Get all pods in the kube-system namespace
	namespace := "kube-system"
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	var cniUsed string

	// Iterate over the containers in each pod
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			// Check if the container image name contains known CNI plugins

			//this returns true when the container image name is in the list of known cni plugins (for further information see the containerImageContainsCNIName() function)
			if containerImageContainsCNIName(container.Image) {

				//the container image that is matched by the list in the function containerImageContainsCNIName()
				cniImage := string(container.Image)

				//runs the switch statement where the strings.Contains returns true
				switch true {
				case strings.Contains(cniImage, "cilium"):
					cniUsed = "Cilium"
				case strings.Contains(cniImage, "calico"):
					cniUsed = "Calico"
				case strings.Contains(cniImage, "weaveworks/weave"):
					cniUsed = "Weave Net"
				case strings.Contains(cniImage, "flannel"):
					cniUsed = "Flannel"
				default:
					cniUsed = "unknown"
				}
			}
		}
	}

	return cniUsed
}

func containerImageContainsCNIName(image string) bool {
	// List of known CNI plugins
	cniPlugins := []string{"calico", "cilium", "flannel", "weave", "kube-router"}

	// Convert the image name to lowercase for case-insensitive matching
	imageLower := strings.ToLower(image)

	// Check if the image name contains any of the known CNI plugins
	for _, plugin := range cniPlugins {

		//returns true if any of the images of containers in the kube-system match the list above
		if strings.Contains(imageLower, plugin) {
			return true

		}
	}

	return false
}

func getPods(clientset *kubernetes.Clientset) *corev1.PodList {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	return pods
}

func getClientSet() *kubernetes.Clientset {

	// create the clientset
	config, err := getKubeconfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func printArt() {

	//initiates clientset
	clientset := getClientSet()

	//gets kubernetes version
	version, err := getKubeVersion(clientset)
	if err != nil {
		panic(err.Error())
	}

	pods := getPods(clientset)

	nodes := getNodeInfo(clientset)

	kubernetesEndpointPort := getKubernetesEndpointPort(clientset)

	isTalos := isTalos(nodes)

	//gets kubernetes distro
	var distro string
	if strings.Contains(version, "k3s") {
		distro = "K3s"
	} else {

		//microk8s usually uses the apiServer port 16443
		if kubernetesEndpointPort == 16443 {
			distro = "MicroK8s"
		} else {
			if isTalos {
				distro = "Talos"
			} else {
				distro = "K8s"
			}
		}

	}

	//gets storage solution used
	storage, err := getStorage(clientset)
	if err != nil {
		panic(err.Error())
	}

	//gets container runtime interface
	containerRuntimeInterface, err := getContainerRuntimeInterface(clientset)
	if err != nil {
		panic(err.Error())
	}

	gitopsTool, err := getGitops(clientset)
	if err != nil {
		panic(err.Error())
	}
	clusterAge, err := getNodeAge(nodes)
	if err != nil {
		panic(err.Error())
	}

	nodeCount, maxPods := getNodeCount(clientset)
	if err != nil {
		panic(err.Error())
	}

	serviceCount := getServiceCount(clientset)
	cniUsed := getCNI(clientset)

	if (cniUsed == "unknown") && (distro == "K3s") {
		//k3s bundles the default flannel cli in a binary in the executable, therefore it's detection via the pod images in the kube-system namespace is impossible
		//therefore, if no known cni is detected in a k3s cluster, it is assumed that the cni is the default bundled flannel
		cniUsed = "Flannel"
	}

	podCount := getPodCount(clientset, pods)
	nonRunningPods := getNonRunningPodCount(clientset, pods)
	namespaceCount := getNamespaceCount(clientset)
	//fetch all values from the various functions above

	var asciiArtColor string

	var red int
	var green int
	var blue int
	switch distro {
	case "MicroK8s":
		red = 233
		green = 84
		blue = 32
	case "K3s":
		red = 255
		green = 198
		blue = 28
	case "Talos":
		red = 249
		green = 42
		blue = 32
	default:
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
		colorCode + asciiArtColor + "    " + "Pod Count: " + resetCode + fmt.Sprint(podCount) + "/" + fmt.Sprint(maxPods),
		colorCode + asciiArtColor + "    " + "Unhealthy pods: " + resetCode + strconv.Itoa(nonRunningPods),
		colorCode + asciiArtColor + "    " + "Namespace Count: " + resetCode + strconv.Itoa(namespaceCount),
		colorCode + asciiArtColor + "    " + "Service Count: " + resetCode + fmt.Sprint(serviceCount),
		colorCode + asciiArtColor + "    " + "Container Runtime Interface: " + resetCode + containerRuntimeInterface,
		colorCode + asciiArtColor + "    " + "Storage: " + resetCode + storage,
		colorCode + asciiArtColor + "    " + "GitOps Tool: " + resetCode + gitopsTool,
		colorCode + asciiArtColor + "    " + "Container Networking Interface: " + resetCode + cniUsed,
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
