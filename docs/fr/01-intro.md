# Construction d'un Runtime de Container à partir de Zéro

Bienvenue dans cet atelier pratique où vous apprendrez à construire les fonctionnalités de base d'un runtime de container à partir de zéro.  
À travers une série d'exercices, vous implémenterez les fonctionnalités essentielles des containers en utilisant les capacités de programmation système de Go.

## Prérequis

- Go 1.23 ou version ultérieure
- Environnement Linux (natif ou via dev container)
- Compréhension de base des concepts de containers
- Accès root/sudo pour les opérations système

## Environnement de Développement

Si vous êtes sur MacOS ou Windows, vous devrez utiliser l'environnement dev container fourni car les exercices nécessitent des fonctionnalités spécifiques à Linux. Deux options sont disponibles :

1. **VS Code / JetBrains DevContainer** : Configuration fournie dans `.devcontainer/`
2. **Docker Compose** : Exécutez `docker compose run --rm -P --build shell` dans le répertoire `.devcontainer/`

## Structure de l'Atelier

L'atelier est divisé en exercices suivants, chacun s'appuyant sur les précédents :

### 1. Gestion des Processus

- [Principes de base de la création de processus](02-process-creation.md)

  - Création de processus parents et enfants
  - Gestion et communication des processus
  - Gestion des erreurs

- [Isolation par Namespaces](03-namespace-isolation.md)
  - Implémentation du namespace PID
  - Namespace UTS pour l'isolation du nom d'hôte
  - Isolation de base des processus

### 2. Fondation du Container

- [Namespaces et Répertoire Racine](04-namespaces-and-chroot.md)

  - Gestion de plusieurs namespaces
  - Implémentation de chroot
  - Configuration de la structure de répertoires

- [Contrôle des Ressources avec cgroups](05-cgroups.md)
  - Limitations CPU
  - Contraintes de mémoire
  - Gestion des ressources des processus

### 3. Fonctionnalités Avancées

- [Gestion des Volumes](06-volumes.md)

  - Implémentation des bind mounts
  - Persistance des volumes
  - Partage de données entre l'hôte et le container

- [Configuration Réseau](07-network.md)
  - Configuration du namespace réseau
  - Paires d'ethernet virtuelles (veth)
  - Capacités réseau de base

## Compilation et Exécution

Commandes de base pour commencer :

```bash
# Compiler le projet
make

# Exécuter un container basique
sudo ./bin/devoxx-docker run alpine /bin/sh
```

## Concepts Clés Abordés

- Isolation et gestion des processus
- Implémentation des namespaces
- Contrôle des ressources avec cgroups
- Opérations sur le système de fichiers
- Configuration réseau
- Gestion des volumes

## Ressources Supplémentaires

- [Linux Namespaces](https://man7.org/linux/man-pages/man7/namespaces.7.html)
- [Control Groups v2](https://www.kernel.org/doc/Documentation/cgroup-v2.txt)
- [Container Networking](https://docs.docker.com/network/)
- [OCI Runtime Specification](https://github.com/opencontainers/runtime-spec)

## Obtenir de l'Aide

- Utilisez `make help` pour voir les commandes disponibles
- Consultez la documentation de chaque exercice
- Référez-vous aux indices et références de commandes dans chaque fichier d'exercice

[Commencer l'atelier](02-process-creation.md)
