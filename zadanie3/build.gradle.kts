plugins {
    application
    kotlin("jvm") version "1.9.23"
    id("com.github.johnrengelman.shadow") version "8.1.1"
}

repositories {
    mavenCentral()
}

dependencies {
    implementation("io.ktor:ktor-server-netty:2.3.10")
    implementation("io.ktor:ktor-server-core:2.3.10")
    implementation("io.ktor:ktor-server-content-negotiation:2.3.10")
    implementation("io.ktor:ktor-serialization-kotlinx-json:2.3.10")

    implementation("dev.kord:kord-core:0.11.0")
    implementation("ch.qos.logback:logback-classic:1.5.3")
}

application {
    mainClass.set("com.Application.ApplicationKt")
}

tasks.shadowJar {
    archiveFileName.set("app.jar")   
    mergeServiceFiles()              
}

tasks.withType<org.jetbrains.kotlin.gradle.tasks.KotlinCompile>().configureEach {
    kotlinOptions.jvmTarget = "17"
}

java.toolchain.languageVersion.set(JavaLanguageVersion.of(17))
