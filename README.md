# imgn!

Welcome to `imgn!` (as in: imagine), an educational app that demonstrates the journey from monolith to containerized microservices to serverless/Function-as-a-Service. The app itself is really simple:

![imgn in action](img/imgn.gif)

It allows you to upload images, viewing them in a gallery. Also, it automatically extracts metadata from the uploaded images. The app is written in Go, using three different ways:

- As a [monolith](monolith/), using TerraForm to deploy into a VM
- As a [containerized microservice](containers/), using Kubernetes
- As a collection of functions, using AWS Lambda