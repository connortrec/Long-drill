package main

import (
	"github.com/pkg/browser"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"strconv"
)

func Show(sync todoist.Sync, c *cli.Context) error {
	item_id, err := strconv.Atoi(c.Args().First())
	if err != nil {
		return err
	}

	idCarrier, err := todoist.SearchByID(sync.Items, item_id)
	item := idCarrier.(todoist.Item)
	if err != nil {
		return err
	}

	colorList := ColorList()
	var projectIds []int
	for _, project := range sync.Projects {
		projectIds = append(projectIds, project.GetID())
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)

	records := [][]string{
		[]string{"ID", IdFormat(item)},
		[]string{"Content", ContentFormat(item)},
		[]string{"Project", ProjectFormat(item.ProjectID, sync.Projects, projectColorHash, c)},
		[]string{"Labels", item.LabelsString(sync.Labels)},
		[]string{"Priority", PriorityFormat(item.Priority)},
		[]string{"DueDate", DueDateFormat(item.DueDateTime(), item.AllDay)},
		[]string{"URL", todoist.GetContentURL(item)},
	}
	defer writer.Flush()

	for _, record := range records {
		writer.Write(record)
	}

	if todoist.HasURL(item) {
		if c.Bool("browse") {
			browser.OpenURL(todoist.GetContentURL(item))
		}
	}
	return nil
}
