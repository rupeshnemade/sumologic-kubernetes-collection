package multilinelogsgenerator

import (
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const bashScriptTemplate = `for i in $(seq 500); do
LONG_STRING="$(cat /dev/urandom | tr -dc ''a-z0-9'' | head -c 30000)";
echo "Dec 13 09:41:08 1st single line...";
echo "Dec 13 09:41:08 2nd single line...";
echo "Dec 14 06:41:08 Exception in thread "main" java.lang.RuntimeException: Something has gone wrong, aborting! ${LONG_STRING} end of the 1st long line
	at com.myproject.module.MyProject.badMethod(MyProject.java:22)
	at com.myproject.module.MyProject.oneMoreMethod(MyProject.java:18)
	at com.myproject.module.MyProject.anotherMethod(MyProject.java:14)
	at com.myproject.module.MyProject.someMethod(MyProject.java:10)";
echo "    at com.myproject.module.MyProject.verylongLine(MyProject.java:100000) ${LONG_STRING} end of the 2nd long line";
echo "    at com.myproject.module.MyProject.main(MyProject.java:6)
Dec 15 09:41:08 another line in loop ${i}";
done`

const image = "busybox"

var terminationGracePeriodSeconds int64 = 3600

const deploymentSleepTime = time.Hour * 24 // how much time we spend sleeping after generating logs in a Deployment

func GetMultilineLogsDeployment(
	namespace string,
	name string,
) appsv1.Deployment {
	var replicas int32 = 1
	appLabels := map[string]string{
		"app": name,
	}
	metadata := metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    appLabels,
	}

	// logsGeneratorAndSleepCommand := fmt.Sprintf("%v;\n sleep %f", bashScriptTemplate, deploymentSleepTime.Seconds())

	podTemplateSpec := corev1.PodTemplateSpec{
		ObjectMeta: metadata,
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
			Containers: []corev1.Container{
				{
					Name:  name,
					Image: image,
					Args:  []string{"/bin/sh", "-c", bashScriptTemplate, ";", "/bin/sh", "-c", "sleep 3600"},
				},
			},
		},
	}
	return appsv1.Deployment{
		ObjectMeta: metadata,
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: appLabels,
			},
			Template: podTemplateSpec,
		},
	}
}
