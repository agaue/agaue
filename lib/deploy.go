package lib

import (
	"fmt"
	// "io/ioutil"
	"log"
	"os/exec"
)

func DeploySite() {
	git, err := exec.LookPath("git")
	if err != nil {
		log.Fatal("git may not installed")
	}
	fmt.Println(git)

	gitInit := exec.Command("git", "init") //TODO: should be fired at init blog, not generate
	gitAdd := exec.Command("git", "add", "-A")
	gitCommit := exec.Command("git", "commit", "-m", "Site Update")
	i, err := gitInit.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	a, err := gitAdd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	c, err := gitCommit.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	//test
	fmt.Println(string(i))
	fmt.Println(string(a))
	fmt.Println(string(c))
}
