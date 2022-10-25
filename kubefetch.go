package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

var kube_url = "https://localhost:38045"

/*
func initClient() {
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}
}
*/

func main() {
	assemblingArt()
}

// gets the amount of nodes in the cluster
func getNodeCount() (nodeCount int) {

	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	resp, err := client.Get(kube_url + "/api/v1/nodes")

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
		log.Println("Connection error while fetching nodes")
	}

	return nodeCount
}

// gets the amount of nodes in the cluster
func getNamespaceCount() (namespaceCount int) {
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	resp, err := client.Get(kube_url + "/api/v1/namespaces")

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
		log.Println("Connection error while fetching namespaces")
	}

	return namespaceCount
}

func getPodCount() (podCount int) {
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	resp, err := client.Get(kube_url + "/api/v1/pods")

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
		log.Println("connection error while fetching pods")
	}

	return podCount
}

func getDistro() (gitVersion string, majorVersion string, minorVersion string, distro string) {

	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	resp, err := client.Get(kube_url + "/version")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	/*
		var gitVersion string
		var majorVersion string
		var minorVersion string
		var distro string
	*/

	if resp.StatusCode == http.StatusOK {

		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		bodyString := string(bodyBytes)
		gitVersionRaw := gjson.Get(bodyString, "gitVersion")

		major_raw := gjson.Get(bodyString, "major")
		minor_raw := gjson.Get(bodyString, "minor")

		if strings.Contains(gitVersionRaw.String(), "k3s") {
			distro = "k3s"
		} else {
			distro = "k8s"
		}

		gitVersion = gitVersionRaw.String()
		majorVersion = major_raw.String()
		minorVersion = minor_raw.String()

	} else {
		log.Println("Connection error while fetching Kubernetes version")
	}
	return gitVersion, majorVersion, minorVersion, distro

}

func getUsedIngress() (ingressUsed string) {
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	respCouldContainTraefik, err := client.Get(kube_url + "/apis/apiextensions.k8s.io/v1/customresourcedefinitions")
	respCouldContainNginx, err := client.Get(kube_url + "/apis/networking.k8s.io/v1/ingressclasses")

	if err != nil {
		log.Fatal(err)
	}

	defer respCouldContainTraefik.Body.Close()
	defer respCouldContainNginx.Body.Close()

	if respCouldContainTraefik.StatusCode == http.StatusOK {

		couldContainTraefik, err := io.ReadAll(respCouldContainTraefik.Body)
		couldContainNginx, err := io.ReadAll(respCouldContainNginx.Body)

		if err != nil {
			log.Fatal(err)
		}

		//managedFields only appears once per node in the json file you get when calling the url, so it is getting used to get the pod count
		var firstAppend int64

		if strings.Contains(string(couldContainTraefik), "traefik.containo.us") {
			ingressUsed = "Traefik"
			firstAppend = 1
		} else {
		}

		if strings.Contains(string(couldContainNginx), "nginx") {
			if firstAppend == 1 {
				ingressUsed = ingressUsed + ", "
			}
			firstAppend = 1
			ingressUsed = ingressUsed + "Nginx"
		} else {
		}

	}
	return ingressUsed

}

func assemblingArt() {

	var gitVersion string
	var majorVersion string
	var minorVersion string
	var distro string
	gitVersion, majorVersion, minorVersion, distro = getDistro()
	var distroArt [17]string
	var asciiArtColor string
	if false == true {
		println(gitVersion)
	}

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
		asciiArtColor = "34"

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

	var nodeCount = getNodeCount()
	var podCount = getPodCount()
	var namespaceCount = getNamespaceCount()
	var ingressUsed = getUsedIngress()

	print("\033[" + asciiArtColor + ";1m" + distroArt[0] + "\033[0m")
	print("\033[" + asciiArtColor + ";1mDistro: \033[0m")
	print(distro)
	print("\n")

	print("\033[" + asciiArtColor + ";1m" + distroArt[1] + "\033[0m")
	print("\033[" + asciiArtColor + ";1mVersion: \033[0m")
	print(majorVersion, ".", minorVersion)
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
	print("\033[" + asciiArtColor + ";1mIngress: \033[0m")
	print(ingressUsed)
	print("\n")

	print("\033[" + asciiArtColor + ";1m" + distroArt[6] + "\033[0m")
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
