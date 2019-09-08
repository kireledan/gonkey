package modules

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kireledan/gonkey/utils"
)

const packageStatePresent = "present"
const packageStateAbsent = "absent"
const packageStateLatest = "latest"

var supportedPackageManagers = []PackageManager{
	BrewManager{},
}

type PackageManager interface {
	IsAvailable() bool
	GetPackages() []string
	Install(pkg string) error
	Remove(pkg string) error
	Update() error
}

type BrewManager struct{}

// IsAvailable checks if homebrew is installed.
func (b BrewManager) IsAvailable() bool {
	_, available := utils.GetBinaryPath("brew")
	return available
}

func (b BrewManager) GetPackages() []string {
	done := make(chan utils.Result)
	utils.ExecCmd("brew list -1", done)

	result := <-done
	packages := result.GetStdout()
	return strings.Split(packages, "\n")
}

func (b BrewManager) Install(pkg string) error {
	done := make(chan utils.Result)
	err := utils.ExecCmd(fmt.Sprintf("brew install %s", pkg), done)

	_ = <-done

	if err != nil {
		return err
	}
	return nil
}

func (b BrewManager) Remove(pkg string) error {
	done := make(chan utils.Result)
	err := utils.ExecCmd(fmt.Sprintf("brew uninstall %s", pkg), done)

	_ = <-done

	if err != nil {
		return err
	}
	return nil
}

func (b BrewManager) Update() error {
	done := make(chan utils.Result)
	err := utils.ExecCmd("brew update", done)

	fmt.Println("waiting... for update status")
	_ = <-done
	fmt.Println("got update status")

	if err != nil {
		return err
	}
	return nil
}

type Package struct {
	name       string
	versionPin string
	state      string
}

func (m Package) execute(args ...string) utils.Result {

	done := make(chan utils.Result)
	go achievePackageState(m.name, m.versionPin, m.state, done)

	result := <-done

	return result
}

func achievePackageState(packageName string, packageVersion string, packageState string, done chan utils.Result) {
	res := new(utils.Result)

	// Detect package manager
	fmt.Println("Detecting package manager...")
	pkgManager, err := detectPackageManager()
	if err != nil {
		res.Stderr = fmt.Sprintf("%s", err)
	}

	fmt.Println("Updating package list...")
	if err := pkgManager.Update(); err != nil {
		res.Stderr = fmt.Sprintf("%s", err)
	}

	fmt.Println("Achieving desired state...")
	// Achieve desired state
	switch packageState {
	case packageStatePresent:
		pkgManager.Install(packageName)
	case packageStateAbsent:
		pkgManager.Remove(packageName)
	case packageStateLatest:
		// TODO
	}

	done <- *res
}

func detectPackageManager() (PackageManager, error) {
	// Get list of all supported managers
	// Check if the package manager exists on the system
	for _, pkg := range supportedPackageManagers {
		manager := pkg
		if manager.IsAvailable() {
			return manager, nil
		}
	}

	return nil, errors.New("No supported package manager available")
}

func (m Package) getModuleName() string {
	return "package"
}

func (m Package) InitModuleFromMap(args map[string]string) Module {
	m.name = args["name"]
	m.versionPin = args["version"]
	m.state = args["state"]
	return m
}
