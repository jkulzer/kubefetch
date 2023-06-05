package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	printArt()
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

func getNodeCount(nodeCount *int) {

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
	*nodeCount = len(nodes.Items)
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

func printArt() {

	//gets kubernetes version
	version, err := getKubeVersion()
	if err != nil {
		panic(err.Error())
	}

	//gets kubernetes distro
	var distro string
	if strings.Contains(version, "k3s") {
		distro = "K3s"
	} else {
		distro = "K8s"
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

	var distroArt [17]string
	var asciiArtColor string
	var podCount int
	var namespaceCount int
	var nodeCount int

	if distro == "K8s" {

		distroArt[0] = "		0MMMNNMMMO               	"
		distroArt[1] = "          .MMMWKko:::cokKWMMM.          	"
		distroArt[2] = "      dMMMN0xl:::::cc:::::lx0NMMMo      	"
		distroArt[3] = "    0MXkdc:::::::::K0:::::::::cdOXMO    	"
		distroArt[4] = "   'MO::::::::::codXXdoc::::::::::0M.   	"
		distroArt[5] = "   OWc:::oXkclONWX0MM0XWXklckXl:::cMk   	"
		distroArt[6] = "  .MO::::::dMMWx::cMM:::kMMMd::::::0M.  	"
		distroArt[7] = "  OMc::::::XMkKMNkkMMkONM0kMX::::::lMx  	"
		distroArt[8] = " .MO::::::xMO::cKMM00MMKc::0Md::::::0M. 	"
		distroArt[9] = " kMc::::::OMX0XWMMWdxWMMWX0XMk::::::lMx 	"
		distroArt[10] = " M0::::kKOOWMkolcXMWWMXclokMWkOKk::::KM 	"
		distroArt[11] = " .No:::::::lNWxcOMKccXMkckWNl:::::::oN. 	"
		distroArt[12] = "   c0c:::::::xXWMWxddxWMWXx:::::::c0:   	"
		distroArt[13] = "     Ox:::::::lWkkO00OkkWl:::::::kk     	"
		distroArt[14] = "       Xo::::cKd::::::::xK:::::oN       	"
		distroArt[15] = "        :0c::::::::::::::::::lK,        	"
		distroArt[16] = "          kklcccccccccccccclOx          	"
		asciiArtColor = "34"

	} else if distro == "K3s" {

		distroArt[0] = " .kWMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMWk.  	"
		distroArt[1] = ":WMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMWc	"
		distroArt[2] = "WMMMMMMMMMMMMMMMMM    MMMMMMMMMMMMMMMMMM	"
		distroArt[3] = "WMMMMMMMMMMMMMMMMX    MMMMMMMMMMMMMMMMMM	"
		distroArt[4] = "WMMMMMMMMMMMMMMMMK    MMMMMMMMMMMMMMMMMM	"
		distroArt[5] = "WMMMMMMMMMMMMMMMMK    MMMMMMMMMMMMMMMMMM	"
		distroArt[6] = "WMMMMMMMMMMMMMMMMK    MMMMMMMMMMMMMMMMMM	"
		distroArt[7] = "WMMMMMMMMMMMMMMMMW    MMMMMMMMMMMMMMMMMM	"
		distroArt[8] = "WMMMMMMMMMMMMMMMMMW00WMMMMMMMMMMMMMMMMMM	"
		distroArt[9] = "WMMMMMMMMMMMMM   lMMMMc   MMMMMMMMMMMMMM	"
		distroArt[10] = "WMMMMMMMMM.       MMMM       .MMMMMMMMMM	"
		distroArt[11] = "WMMMMMx        'lXMMMMKl.        kMMMMMM	"
		distroArt[12] = "WMMMMM     .cOWMMMMMMMMMMWkc.    .MMMMMM	"
		distroArt[13] = "XMMMMMXlcxNMMMMMMMMMMMMMMMMMMNxclXMMMMMN	"
		distroArt[14] = " oMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMo 	"
		distroArt[15] = "    ''''''''''''''''''''''''''''''''    	"
		distroArt[16] = "										 	"
		asciiArtColor = "93"

	} else {

		distroArt[0] = "		0MMMNNMMMO               	"
		distroArt[1] = "          .MMMWKko:::cokKWMMM.          	"
		distroArt[2] = "      dMMMN0xl:::::cc:::::lx0NMMMo      	"
		distroArt[3] = "    0MXkdc:::::::::K0:::::::::cdOXMO    	"
		distroArt[4] = "   'MO::::::::::codXXdoc::::::::::0M.   	"
		distroArt[5] = "   OWc:::oXkclONWX0MM0XWXklckXl:::cMk   	"
		distroArt[6] = "  .MO::::::dMMWx::cMM:::kMMMd::::::0M.  	"
		distroArt[7] = "  OMc::::::XMkKMNkkMMkONM0kMX::::::lMx  	"
		distroArt[8] = " .MO::::::xMO::cKMM00MMKc::0Md::::::0M. 	"
		distroArt[9] = " kMc::::::OMX0XWMMWdxWMMWX0XMk::::::lMx 	"
		distroArt[10] = " M0::::kKOOWMkolcXMWWMXclokMWkOKk::::KM 	"
		distroArt[11] = " .No:::::::lNWxcOMKccXMkckWNl:::::::oN. 	"
		distroArt[12] = "   c0c:::::::xXWMWxddxWMWXx:::::::c0:   	"
		distroArt[13] = "     Ox:::::::lWkkO00OkkWl:::::::kk     	"
		distroArt[14] = "       Xo::::cKd::::::::xK:::::oN       	"
		distroArt[15] = "        :0c::::::::::::::::::lK,        	"
		distroArt[16] = "          kklcccccccccccccclOx          	"
		asciiArtColor = "34"
	}

	getNodeCount(&nodeCount)
	getPodCount(&podCount)
	getNamespaceCount(&namespaceCount)
	//fetch all values from the various functions above

	print("\033[" + asciiArtColor + ";1m" + distroArt[0] + "\033[0m")
	print("\033[" + asciiArtColor + ";1mDistro: \033[0m")
	print(distro)
	print("\n")

	print("\033[" + asciiArtColor + ";1m" + distroArt[1] + "\033[0m")
	print("\033[" + asciiArtColor + ";1mVersion: \033[0m")
	print(version)
	print("\n")

	print("\033[" + asciiArtColor + ";1m" + distroArt[2] + "\033[0m")
	print("\033[" + asciiArtColor + ";1mNode Count: \033[0m")
	print(nodeCount)
	print("\n")

	print("\033[" + asciiArtColor + ";1m" + distroArt[3] + "\033[0m")
	print("\033[" + asciiArtColor + ";1mPod Count: \033[0m")
	print(podCount)
	print("\n")

	print("\033[" + asciiArtColor + ";1m" + distroArt[4] + "\033[0m")
	print("\033[" + asciiArtColor + ";1mNamespace Count: \033[0m")
	print(namespaceCount)
	print("\n")

	print("\033[" + asciiArtColor + ";1m" + distroArt[5] + "\033[0m")
	print("\033[" + asciiArtColor + ";1mContainer Runtime Interface: \033[0m")
	print(containerRuntimeInterface)
	print("\n")

	print("\033[" + asciiArtColor + ";1m" + distroArt[6] + "\033[0m")
	print("\033[" + asciiArtColor + ";1mStorage: \033[0m")
	print(storage)
	print("\n")

	print("\033[" + asciiArtColor + ";1m" + distroArt[7] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[8] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[9] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[10] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[11] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[12] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[13] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[14] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[15] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[16] + "\033[0m")
	print("\n")
	print("")
}
