package main

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

var kube_url = "http://localhost:8001"

func main() {
	getDistro()

}

// gets the amount of nodes in the cluster
func getNodeCount() (nodeCount int) {

	resp, err := http.Get(kube_url + "/api/v1/nodes")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		bodyString := string(bodyBytes)

		//managedFields only appears once per node in the json file you get when calling the url, so it is getting used to get the node count
		nodeCount = strings.Count(bodyString, "managedFields")

	} else {
		print("connection error")
	}

	return nodeCount
}

// gets the amount of nodes in the cluster
func getNamespaceCount() (namespaceCount int) {

	resp, err := http.Get(kube_url + "/api/v1/namespaces")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		bodyString := string(bodyBytes)

		//managedFields only appears once per node in the json file you get when calling the url, so it is getting used to get the namespace count
		namespaceCount = strings.Count(bodyString, "managedFields")

	} else {
		print("connection error")
	}

	return namespaceCount
}
func getPodCount() (podCount int) {

	resp, err := http.Get(kube_url + "/api/v1/pods")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		bodyString := string(bodyBytes)

		//managedFields only appears once per node in the json file you get when calling the url, so it is getting used to get the pod count
		podCount = strings.Count(bodyString, "managedFields")

	} else {
		print("connection error")
	}

	return podCount
}

func getDistro() {

	resp, err := http.Get(kube_url + "/version")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		bodyString := string(bodyBytes)
		kube_version := gjson.Get(bodyString, "gitVersion")

		major_version := gjson.Get(bodyString, "major")
		minor_version := gjson.Get(bodyString, "minor")

		var distro string

		if strings.Contains(kube_version.String(), "k3s") {
			distro = "k3s"
		} else {
			distro = "k8s"
		}

		assemblingArt(distro, kube_version.String(), major_version.String(), minor_version.String())
	} else {
		print("connection error")
	}

}

func assemblingArt(distro string, kube_version string, major_version string, minor_version string) {
	var distroArt [17]string

	if distro == "k8s" {

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

	} else if distro == "k3s" {

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
	}

	var nodeCount = getNodeCount()
	var podCount = getPodCount()
	var namespaceCount = getNamespaceCount()

	print(distroArt[0])
	print("Distro: ")
	print(distro)
	print("\n")

	print(distroArt[1])
	print("Version: ")
	print(major_version, ".", minor_version)
	print("\n")

	print(distroArt[2])
	print("Node Count: ")
	print(nodeCount)
	print("\n")

	print(distroArt[3])
	print("Pod Count: ")
	print(podCount)
	print("\n")

	print(distroArt[4])
	print("Namespace Count: ")
	print(namespaceCount)
	print("\n")
	print(distroArt[5])
	print("\n")
	print(distroArt[6])
	print("\n")
	print(distroArt[7])
	print("\n")
	print(distroArt[8])
	print("\n")
	print(distroArt[9])
	print("\n")
	print(distroArt[10])
	print("\n")
	print(distroArt[11])
	print("\n")
	print(distroArt[12])
	print("\n")
	print(distroArt[13])
	print("\n")
	print(distroArt[14])
	print("\n")
	print(distroArt[15])
	print("\n")
	print(distroArt[16])
	print("\n")
	print("")
}
