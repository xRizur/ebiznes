import play.sbt.PlayImport._

name := "play-scala-app"

version := "1.0.0"

scalaVersion := "2.13.6"

libraryDependencies ++= Seq(
  guice,
  "com.typesafe.play" %% "play" % "2.8.15",
  "com.typesafe.play" %% "play-json" % "2.9.2"
)

enablePlugins(PlayScala)
