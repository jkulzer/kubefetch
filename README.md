<br />
<div align="center">
    <img src="https://cncf-branding.netlify.app/img/projects/kubernetes/icon/color/kubernetes-icon-color.svg" alt="Logo" width="80" height="80">
  </a>

<h1 align="center">Kubefetch</h1>

  <p align="center">
    A neofetch-inspired application that displays stats about your Kubernetes cluster.
    <br />
    ·
    <a href="https://github.com/jkulzer/kubefetch/issues">Report Bug</a>
    ·
    <a href="https://github.com/jkulzer/kubefetch/issues">Request Feature</a>
  </p>
</div>


## Feel free to open an issue if your Kubernetes distro or some other aspect does not get recognized!

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

### Uses:

Neofetch for Kubernetes Clusters

<br>

![](https://github.com/jkulzer/kubefetch/blob/main/kubefetch.png?raw=true)

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

### Installation

1. With `make`
    Clone the repo
    ```sh
    git clone https://github.com/jkulzer/kubefetch
    cd kubefetch
    ```

2. Build and install the binary
    ```sh
    make build
    sudo make install
    ```

3. To uninstall, run:
    ```
    sudo make uninstall
    make clean
    ```

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- ROADMAP -->
## Roadmap

- [x] Auth via kubeconfig
- [ ] AUR package
- [ ] More displayed info
	- [ ] CNI used
	- [x] CRI used
	- [x] Storage Solution used

See the [open issues](https://github.com/jkulzer/kubefetch/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#top">back to top</a>)</p>


<!-- LICENSE -->
## License

Distributed under the GNU GPL v3 License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Project Link: [https://github.com/jkulzer/kubefetch](https://github.com/jkulzer/kubefetch)

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* [k8s@home](https://k8s-at-home.com/)

<p align="right">(<a href="#top">back to top</a>)</p>
