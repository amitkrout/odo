package e2e

import (
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const javaFiles = "examples/binary/java/"

var _ = Describe("odoJavaE2e", func() {
	const t = "java"
	var projName = fmt.Sprintf("odo-%s", t)

	// contains a minimal javaee app
	const warGitRepo = "https://github.com/lordofthejars/book-insultapp"

	// contains a minimal javalin app
	const jarGitRepo = "https://github.com/geoand/javalin-helloworld"

	// Create a separate project for Java
	Context("create java project", func() {
		It("should create a new java project", func() {
			session := runCmdShouldPass("odo project create " + projName)
			Expect(session).To(ContainSubstring(projName))
			waitForCmdOut("odo project set "+projName, 4, false, func(output string) bool {
				return strings.Contains(output, "Already on project : "+projName)
			})
		})
	})

	// Test Java
	Context("odo component creation", func() {
		It("Should be able to deploy a git repo that contains a wildfly application", func() {

			// Deploy the git repo / wildfly example
			cmpCreateLog := runCmdShouldPass("odo create wildfly javaee-git-test --git " + warGitRepo + " -w")
			Expect(cmpCreateLog).ShouldNot(ContainSubstring("This may take few moments to be ready"))
			buildName := getBuildNameUsingOc("javaee-git-test")
			Expect(buildName).To(ContainSubstring("javaee-git-test"))
			buildStatus := runCmdShouldPass("oc get build " + buildName)
			Expect(buildStatus).To(ContainSubstring("Complete"))

			cmpList := runCmdShouldPass("odo list")
			Expect(cmpList).To(ContainSubstring("javaee-git-test"))

			// Push changes
			runCmdShouldPass("odo push")

			// Create a URL
			runCmdShouldPass("odo url create")
			routeURL := determineRouteURL()

			// Ping said URL
			responsePing := retryingForOutputMatchStringOfHTTPResponse(routeURL, "Insult", 30, 1)
			Expect(responsePing).Should(BeTrue())

			// Delete the component
			runCmdShouldPass("odo delete javaee-git-test -f")
		})

		It("Should be able to deploy a git repo that contains a wildfly application without wait flag", func() {
			// Deploy the git repo / wildfly example
			cmpCreateLog := runCmdShouldPass("odo create wildfly wo-wait-javaee-git-test --git " + warGitRepo)
			Expect(cmpCreateLog).Should(ContainSubstring("This may take few moments to be ready"))
			buildName := getBuildNameUsingOc("wo-wait-javaee-git-test")
			Expect(buildName).To(ContainSubstring("wo-wait-javaee-git-test"))

			buildStatus := getBuildParametersValueUsingOc("wo-wait-javaee-git-test")
			Expect(buildStatus).To(ContainSubstring("Pending"))

			dcName := getDcNameUsingOc("wo-wait-javaee-git-test")
			// For waiting until the deployment starts
			for {
				time.Sleep(5 * time.Second)
				dcStatus := getDcStatusValueUsingOc("wo-wait-javaee-git-test")
				if dcStatus == "1" {
					break
				}
			}
			// following the logs and waiting for the build to finish
			runCmdShouldPass("oc logs --version=1 dc/" + dcName)

			cmpList := runCmdShouldPass("odo list")
			Expect(cmpList).To(ContainSubstring("wo-wait-javaee-git-test"))

			// Push changes
			runCmdShouldPass("odo push")

			// Create a URL
			runCmdShouldPass("odo url create")
			routeURL := determineRouteURL()

			// Ping said URL
			responsePing := retryingForOutputMatchStringOfHTTPResponse(routeURL, "Insult", 30, 1)
			Expect(responsePing).Should(BeTrue())

			// Delete the component
			runCmdShouldPass("odo delete wo-wait-javaee-git-test -f")
		})

		It("Should be able to deploy a .war file using wildfly", func() {
			cmpCreateLog := runCmdShouldPass("odo create wildfly javaee-war-test --binary " + javaFiles + "/wildfly/ROOT.war -w")
			Expect(cmpCreateLog).ShouldNot(ContainSubstring("This may take few moments to be ready"))
			dcName := getDcNameUsingOc("javaee-war-test")
			// Following the logs
			runCmdShouldPass("oc logs --version=1 dc/" + dcName)

			cmpList := runCmdShouldPass("odo list")
			Expect(cmpList).To(ContainSubstring("javaee-war-test"))

			// Push changes
			runCmdShouldPass("odo push")

			// Create a URL
			runCmdShouldPass("odo url create")
			routeURL := determineRouteURL()

			// Ping said URL
			responsePing := retryingForOutputMatchStringOfHTTPResponse(routeURL, "Sample", 30, 1)
			Expect(responsePing).Should(BeTrue())

			// Delete the component
			runCmdShouldPass("odo delete javaee-war-test -f")
		})

		It("Should be able to deploy a git repo that contains a java uberjar application using openjdk", func() {
			importOpenJDKImage()

			// Deploy the git repo / wildfly example
			cmpCreateLog := runCmdShouldPass("odo create openjdk18 uberjar-git-test --git " + jarGitRepo + " -w")
			Expect(cmpCreateLog).ShouldNot(ContainSubstring("This may take few moments to be ready"))
			buildName := getBuildNameUsingOc("uberjar-git-test")
			Expect(buildName).To(ContainSubstring("uberjar-git-test"))
			buildStatus := runCmdShouldPass("oc get build " + buildName)
			Expect(buildStatus).To(ContainSubstring("Complete"))

			cmpList := runCmdShouldPass("odo list")
			Expect(cmpList).To(ContainSubstring("uberjar-git-test"))

			// Push changes
			runCmdShouldPass("odo push")

			// Create a URL
			runCmdShouldPass("odo url create --port 8080")
			routeURL := determineRouteURL()

			// Ping said URL
			responsePing := retryingForOutputMatchStringOfHTTPResponse(routeURL, "Hello World", 30, 1)
			Expect(responsePing).Should(BeTrue())

			// Delete the component
			runCmdShouldPass("odo delete uberjar-git-test -f")
		})

		It("Should be able to deploy a spring boot uberjar file using openjdk", func() {
			importOpenJDKImage()

			cmpCreateLog := runCmdShouldPass("odo create openjdk18 sb-jar-test --binary " + javaFiles + "/openjdk/sb.jar -w")
			Expect(cmpCreateLog).ShouldNot(ContainSubstring("This may take few moments to be ready"))

			dcName := getDcNameUsingOc("sb-jar-test")
			Expect(dcName).To(ContainSubstring("sb-jar-test"))

			cmpList := runCmdShouldPass("odo list")
			Expect(cmpList).To(ContainSubstring("sb-jar-test"))

			// Push changes
			runCmdShouldPass("odo push")

			// Create a URL
			runCmdShouldPass("odo url create --port 8080")
			routeURL := determineRouteURL()

			// Ping said URL
			responsePing := retryingForOutputMatchStringOfHTTPResponse(routeURL, "HTTP Booster", 30, 1)
			Expect(responsePing).Should(BeTrue())

			// Delete the component
			runCmdShouldPass("odo delete sb-jar-test -f")
		})

	})

	// Delete the project
	Context("java project delete", func() {
		It("should delete java project", func() {
			odoDeleteProject(projName)
		})
	})
})

func importOpenJDKImage() {
	// we need to import the openjdk image which is used for jars because it's not available by default
	runCmdShouldPass("oc import-image openjdk18 --from=registry.access.redhat.com/redhat-openjdk-18/openjdk18-openshift:1.5 --confirm")
	runCmdShouldPass("oc annotate istag/openjdk18:latest tags=builder --overwrite")
}
