package run

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/executil"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/template"

	"github.com/wangweihong/eazycloud/apis/ikubeagent"
)

var (
	templateSuffixToWrite     = []string{".template", ".yaml", ".yml", ".sh"}
	fileSuffixToApply         = []string{".yaml", ".yml"}
	fileSuffixToWrite         = []string{".yaml", ".yml", ".sh", ".config"}
	templateSuffixToExecute   = ".template"
	kustomizeTemplateNameList = []string{"kustomization.yaml.template", "kustomization.yaml"}
)

func isFileSuffixToWrite(fileName string) bool {
	for _, v := range fileSuffixToWrite {
		if strings.HasSuffix(fileName, v) {
			return true
		}
	}
	return false
}
func isFileSuffixToApply(fileName string) bool {
	for _, v := range fileSuffixToApply {
		if strings.HasSuffix(fileName, v) {
			return true
		}
	}
	return false
}

func TrimError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Errorf("%v", strings.Replace(err.Error(), "\n", " ", -1))
}

func RunTemplateError(err error, pi template.ProcessorInterface, ctx map[string]any) error {
	if err != nil {
		return err
	}

	switch p := pi.(type) {
	case template.FileProcessor:
		err = parseRunTemplateFile(&p, ctx)
		if err != nil {
			return fmt.Errorf("parse template %v err:%v", p.Name(), err)
		}
	case *template.FileProcessor:
		err = parseRunTemplateFile(p, ctx)
		if err != nil {
			return fmt.Errorf("parse template %v err:%v", p.Name(), err)
		}
	case *template.DirectoryProcessor:
		err = parseRunTemplateDir(p, ctx)
		if err != nil {
			return fmt.Errorf("parse template dir %v err:%v", p.Name(), err)
		}
	case template.DirectoryProcessor:
		err = parseRunTemplateDir(&p, ctx)
		if err != nil {
			return fmt.Errorf("parse template dir %v err:%v", p.Name(), err)
		}
	}
	return nil
}

func writeTemplateFile(tf *template.FileProcessor, ctx map[string]any, fileMode os.FileMode) error {
	log.Infof("generate template %v under path %v", tf.Name(), tf.FilePath)

	tp := tf.SetContexts(ctx)
	log.Infof("template context %v", tp.Context)

	if err := tp.SetFileMode(fileMode).LocateToDisk().Error(); err != nil {
		log.Errorf("generate template %v error: %v", tp.Name(), err)
		return err
	}
	return nil
}

func WriteTemplate(p any, ctx map[string]any, fileMode os.FileMode) error {
	switch pt := p.(type) {
	case *template.FileProcessor:
		return writeTemplateFile(pt, ctx, fileMode)
	case template.FileProcessor:
		return writeTemplateFile(&pt, ctx, fileMode)
	case *template.DirectoryProcessor:
		return parseRunTemplateDir(pt, ctx)
	case template.DirectoryProcessor:
		return parseRunTemplateDir(&pt, ctx)
	}
	return fmt.Errorf("invaild template type")
}

func RunTemplate(p any, ctx map[string]any) error {
	switch pt := p.(type) {
	case *template.FileProcessor:
		return parseRunTemplateFile(pt, ctx)
	case template.FileProcessor:
		return parseRunTemplateFile(&pt, ctx)
	case *template.DirectoryProcessor:
		return parseRunTemplateDir(pt, ctx)
	case template.DirectoryProcessor:
		return parseRunTemplateDir(&pt, ctx)
	}
	return fmt.Errorf("invaild template type")
}

func parseRunTemplateFile(p *template.FileProcessor, ctx map[string]any) error {
	if p == nil {
		return errors.Errorf("template processor is nil")
	}

	p = p.SetContexts(ctx).LocateToDisk()

	// ignore non k8s resource file
	if !isFileSuffixToApply(p.FilePath) {
		return nil
	}
	kustomizePath := filepath.Join(filepath.Dir(p.FilePath), "kustomization.yaml")
	if _, err := os.Stat(kustomizePath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("stat kustomizePath:%v fail:%v", kustomizePath, err)
		}
		args := []string{"--kubeconfig", ikubeagent.KubeConfigPath, "apply", "-f", p.FilePath}
		if _, _, err := executil.ExecuteCmdSplitStdoutStderr(ikubeagent.KubectlBinary, args, 0); err != nil {
			log.Errorf("run command [%v:%v] fail:%v", ikubeagent.KubectlBinary, args, TrimError(err))
			return TrimError(err)
		}
	} else {
		bashCommand := fmt.Sprintf("%v build %v | %v --kubeconfig %v apply -f -",
			ikubeagent.KustomizeBinary, p.FilePath, ikubeagent.KubectlBinary, ikubeagent.KubeConfigPath)
		args := []string{"-c", bashCommand}
		if _, _, err := executil.ExecuteCmdSplitStdoutStderr("/bin/bash", args, 0); err != nil {
			log.Errorf("run command [%v:%v] fail:%v", "/bin/bash", args, TrimError(err))
			return TrimError(err)
		}
	}
	return nil
}

func parseRunTemplateDir(p *template.DirectoryProcessor, ctx map[string]any) error {
	if p == nil {
		return errors.Errorf("template processor is nil")
	}

	p = p.SetContexts(ctx).LocateToDisk()

	kustomizePath := filepath.Join(p.LocateParsedDir, "kustomization.yaml")
	if _, err := os.Stat(kustomizePath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("stat kustomizePath:%v fail:%v", kustomizePath, err)
		}
		args := []string{"--kubeconfig", ikubeagent.KubeConfigPath, "apply", "-f", p.LocateParsedDir}
		if _, stderr, err := executil.ExecuteCmdSplitStdoutStderr(ikubeagent.KubectlBinary, args, 120); err != nil {
			log.Errorf("run command [%v:%v] fail:%v,stderr,%v", ikubeagent.KubectlBinary, args, TrimError(err), stderr)
			return TrimError(err)
		}

	} else {
		bashCommand := fmt.Sprintf("%v build %v | %v --kubeconfig %v apply  -f -", "kustomize",
			p.LocateParsedDir, ikubeagent.KubectlBinary, ikubeagent.KubeConfigPath)
		args := []string{"-c", bashCommand}
		if _, stderr, err := executil.ExecuteCmdSplitStdoutStderr("/bin/bash", args, 0); err != nil {
			log.Errorf("run command [%v:%v] fail:%v,stderr:%v", "/bin/bash", args, TrimError(err), stderr)
			return TrimError(err)
		}
	}
	return nil
}
