package version130

import (
	"os"

	"github.com/wangweihong/gotoolbox/pkg/executil"
	"github.com/wangweihong/gotoolbox/pkg/log"

	"github.com/wangweihong/eazycloud/apis/ikubeagent"
	"github.com/wangweihong/eazycloud/internal/pkg/run"
)

func recordKubeadmLog(stdout string, stderr string) {
	logdata := stdout + "\n" + stderr
	if err := os.WriteFile(ikubeagent.KubeadmLogPath, []byte(logdata), 0755); err != nil {
		log.Errorf("save kubeadm join log to path :%v fail:%v", ikubeagent.KubeadmLogPath, err)
	}
}

func runKubeadmJoin(joinCmd string, isControlPlane bool, nodeName string) error {
	var stderr string
	var stdout string
	var err error

	// make dir manifests, so we can ignore kubelet endless not found manifests error
	_ = os.MkdirAll(ikubeagent.ManifestsDir, 0755)

	// control plane join cluster need to put keepalived/haproxy yaml to manifests first
	if isControlPlane {
		joinCmd += " --ignore-preflight-errors DirAvailable--etc-kubernetes-manifests"
	}
	joinCmd += " --node-name " + nodeName

	args := []string{"-c"}
	args = append(args, joinCmd)

	stdout, stderr, err = executil.ExecuteCmdSplitStdoutStderr("bash", args, 600)
	if err != nil {
		recordKubeadmLog(stdout, stderr)
		log.Errorf("join cluster fail:%v,read %v file for details", run.TrimError(err), ikubeagent.KubeadmLogPath)
		return run.TrimError(err)
	}
	recordKubeadmLog(stdout, stderr)
	log.Infof("run  work for kubeadm join cluster complete")
	return nil
}
