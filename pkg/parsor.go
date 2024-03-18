package pkg

import (
	// yaml

	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"

	"metatask/pkg/schema"
)

type ParsorConfig struct{}

type ParsorInterface interface {
	ParsefromMetaTaskFile(r *schema.FileRoot, c *ParsorConfig) error
}

type V1YamlParsor struct {
	l *logrus.Logger
}

func NewV1YamlParsor(
	l *logrus.Logger,
) *V1YamlParsor {
	return &V1YamlParsor{
		l: l,
	}
}

func (p *V1YamlParsor) Parse(reader io.Reader, tr *schema.TreeRoot, fr *schema.FileRoot, c *ParsorConfig) error {
	// check if the file exists
	// if it does, return an error
	// if it fails, return an error
	//err = json.NewDecoder(fp).Decode(&m)
	err := yaml.NewDecoder(reader).Decode(&fr)
	if err != nil {
		p.l.Error("Error decoding file: ", err)
		return err
	}
	// convert the file root to a tree root
	p.l.Debug("Converting file root to tree root")
	err = p.ParsefromMetaTaskRoot(fr, tr, c)
	if err != nil {
		return err
	}

	return nil
}

func (p *V1YamlParsor) ParsefromMetaTaskRoot(fileRoot *schema.FileRoot, treeRoot *schema.TreeRoot, c *ParsorConfig) error {
	taskMap := make(map[string]*schema.TreeTask)
	var tasksWithDeps []*schema.TreeTask // To handle dependencies after all tasks are created

	// First, create all tasks without dependencies
	for _, ft := range fileRoot.Tasks {
		tt := &schema.TreeTask{
			Name:        ft.Name,
			Command:     ft.Command,
			Description: ft.Description,
		}
		taskMap[ft.Name] = tt
		treeRoot.Tasks = append(treeRoot.Tasks, tt)
	}

	// Now, resolve dependencies
	for _, ft := range fileRoot.Tasks {
		tt := taskMap[ft.Name]
		for _, depName := range ft.DependsOn {
			p.l.Debug("TaskName: ", ft.Name)
			p.l.Debug("DepName: ", depName)
			depTask, exists := taskMap[depName]
			if !exists {
				errMsg := fmt.Sprintf("dependency %s not found for task %s", depName, ft.Name)
				p.l.Error(errMsg)
				return fmt.Errorf(errMsg)
			}
			tt.DependsOn = append(tt.DependsOn, depTask)
		}
		if len(ft.DependsOn) > 0 {
			tasksWithDeps = append(tasksWithDeps, tt)
		}
	}

	// Optional: Check for circular dependencies
	if err := checkForCircularDependencies(tasksWithDeps); err != nil {
		p.l.Error("Circular dependency detected: ", err)
		return err
	}

	for _, task := range treeRoot.Tasks {
		task.DependsOn = taskMap[task.Name].DependsOn
	}

	return nil
}

func checkForCircularDependencies(tasks []*schema.TreeTask) error {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var dfs func(*schema.TreeTask) bool
	dfs = func(task *schema.TreeTask) bool {
		visited[task.Name] = true
		recStack[task.Name] = true

		for _, dep := range task.DependsOn {
			if !visited[dep.Name] {
				if dfs(dep) {
					return true // If a cycle is detected in a deeper call
				}
			} else if recStack[dep.Name] {
				return true // Cycle detected
			}
		}

		recStack[task.Name] = false
		return false
	}

	for _, task := range tasks {
		if !visited[task.Name] {
			if dfs(task) {
				return fmt.Errorf("a cycle detected involving task %s", task.Name)
			}
		}
	}

	return nil
}