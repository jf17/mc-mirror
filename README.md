# mc-mirror

`mc-mirror` is a simple HTTP proxy for mirroring Maven Central with local artifact caching. It intercepts `GET` and `HEAD` requests from Maven clients, caches the downloaded files, and serves them directly on subsequent requests.

## 🧩 Features

* ✅ Supports `GET` and `HEAD` requests
* 🗂 Local artifact caching in the `./storage` directory
* 🔌 Easy port configuration via `-port` flag
* ⚡ Works as a "lazy" mirror — downloads files on demand
* 🐳 **Planned: Docker Compose support for easier deployment**

## 🚀 Installation and Run

```bash
git clone https://github.com/jf17/mc-mirror.git
cd mc-mirror
go run mc-mirror.go -port=8080
```

### Command-Line Arguments

| Parameter | Description                   | Default Value |
| --------- | ----------------------------- | ------------- |
| `-port`   | Port on which the server runs | `8080`        |

## 🔧 Maven Configuration

To use `mc-mirror` as a mirror, add the following block to your `~/.m2/settings.xml` file:

```xml
<settings>
  <mirrors>
    <mirror>
      <id>mc-mirror</id>
      <url>http://localhost:8080</url>
      <mirrorOf>*</mirrorOf>
    </mirror>
  </mirrors>
</settings>
```

## 🔧 Gradle Configuration

To use `mc-mirror` with Gradle, modify your `repositories` block in `build.gradle` (Groovy DSL):

```groovy
repositories {
    maven {
        url "http://localhost:8080"
    }
}
```

Or in `build.gradle.kts` (Kotlin DSL):

```kotlin
repositories {
    maven {
        url = uri("http://localhost:8080")
    }
}
```

> 💡 Note: If you're also using other repositories, make sure to list them after `mc-mirror`, since Gradle checks repositories in order.


## 💡 Example Usage

Once the server is running and your build tool is configured, you can run commands like:

```bash
mvn dependency:resolve
```

or for Gradle:

```bash
./gradlew build
```

The first time, dependencies will be fetched from Maven Central and cached. On subsequent requests, they will be served from the local `./storage` directory.
