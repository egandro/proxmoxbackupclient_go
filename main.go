package main

import (
	"flag"
	"fmt"
	"os"
	"proxmoxbackupgo/pkg/pbs"
	"runtime"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/tawesoft/golib/v2/dialog"
)

func main() {
	// Define command-line flags
	baseURLFlag := flag.String("baseurl", "", "Base URL for the proxmox backup server, example: https://192.168.1.10:8007")
	certFingerprintFlag := flag.String("certfingerprint", "", "Certificate fingerprint for SSL connection, example: ea:7d:06:f9...")
	authIDFlag := flag.String("authid", "", "Authentication ID (PBS Api token)")
	secretFlag := flag.String("secret", "", "Secret for authentication")
	datastoreFlag := flag.String("datastore", "", "Datastore name")
	namespaceFlag := flag.String("namespace", "", "Namespace (optional)")
	backupSourceDirFlag := flag.String("backupdir", "", "Backup source directory, must not be symlink")
	pxarOut := flag.String("pxarout", "", "Output PXAR archive for debug purposes (optional)")

	// Parse command-line flags
	flag.Parse()

	// Validate required flags
	if *baseURLFlag == "" || *certFingerprintFlag == "" || *authIDFlag == "" || *secretFlag == "" || *datastoreFlag == "" || *backupSourceDirFlag == "" {

		if runtime.GOOS == "windows" {
			usage := "All options are mandatory:\n"
			flag.VisitAll(func(f *flag.Flag) {
				usage += "-" + f.Name + " " + f.Usage + "\n"
			})
			dialog.Error(usage)
		} else {
			fmt.Println("All options are mandatory")

			flag.PrintDefaults()
		}
		os.Exit(1)
	}

	if runtime.GOOS == "windows" {

		go systray.Run(func() {
			systray.SetIcon(ICON)
			systray.SetTooltip("PBSGO Backup running")
			beeep.Notify("Proxmox Backup Go", fmt.Sprintf("Backup started"), "")
		},
			func() {

			})
	}

	err, newchunkLoad, reusechunkLoad := pbs.RunPbsBackupClient(baseURLFlag, certFingerprintFlag, authIDFlag, secretFlag, datastoreFlag, namespaceFlag, backupSourceDirFlag, pxarOut)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	fmt.Printf("New %d , Reused %d\n", newchunkLoad, reusechunkLoad)
	if runtime.GOOS == "windows" {
		systray.Quit()
		beeep.Notify("Proxmox Backup Go", fmt.Sprintf("Backup complete\nChunks New %d , Reused %d\n", newchunkLoad, reusechunkLoad), "")
	}

}
