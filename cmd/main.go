package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"gopkg.in/yaml.v2"
)

type Params struct {
	Params map[string]string `yaml:"params"`
}

func DeleteParam(svc *ssm.SSM, name string) error {
	_, err := svc.DeleteParameter(
		&ssm.DeleteParameterInput{
			Name: aws.String(name),
		},
	)
	return err
}

func WriteParam(svc *ssm.SSM, name string, value string) error {
	_, err := svc.PutParameter(
		&ssm.PutParameterInput{
			Name:      aws.String(name),
			Value:     aws.String(value),
			Type:      aws.String("String"),
			Overwrite: aws.Bool(true),
		},
	)
	return err
}

// func GetParamHistory(svc *ssm.SSM, )

func main() {
	pathFlag := flag.String("path", "params.yaml", "Path to parameter yaml file")
	flag.Parse()

	absPath, err := filepath.Abs(*pathFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	paramsData, err := ioutil.ReadFile(absPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	params := Params{}
	yaml.Unmarshal(paramsData, &params)

	sess, err := session.NewSession()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ssmSvc := ssm.New(sess)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for k, v := range params.Params {
		fmt.Printf("[%s]: %s\n", k, v)
		err := WriteParam(ssmSvc, k, v)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	out, err := ssmSvc.DescribeParameters(&ssm.DescribeParametersInput{
		ParameterFilters: []*ssm.ParameterStringFilter{
			&ssm.ParameterStringFilter{
				Key:    aws.String("Name"),
				Option: aws.String("BeginsWith"),
				Values: []*string{
					aws.String("/ci/api/"),
					aws.String("/qa/api/"),
				},
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hist, err := ssmSvc.GetParameterHistory(&ssm.GetParameterHistoryInput{
		Name: aws.String("/ci/api/DB_PASSWORD"),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(out)
	fmt.Println(hist)
	fmt.Println(params.Params)
	os.Exit(0)
}
