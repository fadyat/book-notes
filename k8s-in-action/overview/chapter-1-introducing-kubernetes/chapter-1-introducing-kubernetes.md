# Chapter 1. Introducing Kubernetes

- Migrating from monoliths to microservices.
- Bigger number of components make it harder to configure and manage them.
- Kubernetes allows developers to deploy there apps as often as they want, without
  requiring any assistance from the ops team.
- Helps to monitoring and rescheduling apps in case of failures.
- It's becoming a standard way of running distributed apps in the cloud.

## 1.1 Understanding the need for a system like Kubernetes

### 1.1.1 Moving from monolithic apps to microservices

- If any part of a monolithic application isn’t scalable, the whole application becomes unscalable, unless you can split
  up the monolith somehow.
- Simple communication between components.
- Separate component development.
- Easier to scale.
- Growing number of components makes it harder to make deployment-related decisions.
- Microservices make harder to debug and trace execution calls. (Zipkin, Jaeger, etc.)
- Opportunity to have different versions of the same libraries, tools used for
  development, etc.

### 1.1.2 Providing a consistent environment to applications

- 



