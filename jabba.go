package jabba

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"text/template"
)

type File struct {
	Path     string
	Template string
	Perm     os.FileMode
}

type User struct {
	Username string
	Key      string
	Sudo     bool
}

func RunOrDie(c string, args ...string) {
	err := Run(c, args...)
	if err != nil {
		log.Fatal(err)
	}
}

func Run(c string, args ...string) error {
	cmd := exec.Command(c, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func WriteFile(file File) {
	f, err := os.Create(file.Path)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("out").Parse(file.Template)
	if err != nil {
		log.Fatal(err)
	}

	out := f
	testing, _ := strconv.ParseBool(os.Getenv("TEST"))
	if testing {
		out = os.Stdout
	}
	err = tmpl.Execute(out, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func AddUser(user User) {
	err := Run("getent", "passwd", user.Username)
	if err == nil {
		fmt.Println("User", user.Username, "already exists")
		return
	}

	Run("adduser", "--gecos", `""`, "--disabled-password", user.Username)

	if user.Sudo {
		Run("usermod", "--append", "--groups", "sudo", user.Username)
	}

	err = os.MkdirAll("/home/"+user.Username+"/.ssh", 0755)
	if err != nil {
		log.Fatal(err)
	}

	WriteFile(File{
		Path:     "/home/" + user.Username + "/.ssh/authorized_keys",
		Template: user.Key,
		Perm:     0644,
	})
}
