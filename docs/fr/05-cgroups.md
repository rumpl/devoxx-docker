# Découverte de la Configuration des cgroups

## Objectif

Dans cet exercice, vous apprendrez à configurer les cgroups pour limiter l'utilisation de la mémoire et du CPU pour un processus. Vous vous concentrerez sur la configuration de `memory.max`, `cpu.max`, et l'ajout du processus à `cgroup.procs`.

## Étapes

### Étape 1 : Configuration des cgroups

1. Créer un nouveau répertoire pour le cgroup :
   - Ajoutez le code suivant pour créer le répertoire du cgroup :
     ```go
     cgroupPath := "/sys/fs/cgroup/devoxx-container"
     if err := os.Mkdir(cgroupPath, 0755); err != nil {
         log.Fatalf("Failed to create cgroup directory: %v", err)
     }
     ```

### Étape 2 : Configurer la Limite de Mémoire

1. Définir la limite de mémoire à 100MB :
   - Ajoutez le code suivant pour définir la limite de mémoire :
     ```go
     if err := os.WriteFile(cgroupPath+"/memory.max", []byte("104857600"), 0644); err != nil {
         log.Fatalf("Failed to set memory limit: %v", err)
     }
     ```

### Étape 3 : Configurer la Limite de CPU

1. Définir la limite de CPU à 50ms par 100ms :
   - Ajoutez le code suivant pour définir la limite de CPU :
     ```go
     if err := os.WriteFile(cgroupPath+"/cpu.max", []byte("50000 100000"), 0644); err != nil {
         log.Fatalf("Failed to set CPU limit: %v", err)
     }
     ```

### Étape 4 : Ajouter le Processus au cgroup

1. Ajouter le processus au cgroup :
   - Ajoutez le code suivant pour ajouter le processus au cgroup :
     ```go
     if err := os.WriteFile(cgroupPath+"/cgroup.procs", []byte(strconv.Itoa(pid)), 0644); err != nil {
         log.Fatalf("Failed to add process to cgroup: %v", err)
     }
     ```

### Étape 5 : Exécuter le Programme

1. Compiler et exécuter le programme Go.

   - Enregistrez le code dans un fichier nommé `main.go`.
   - Ouvrez un terminal et naviguez vers le répertoire contenant `main.go`.
   - Exécutez les commandes suivantes pour compiler et exécuter le programme :
     ```sh
     go build -o devoxx-cgroup
     sudo ./devoxx-cgroup
     ```

2. Observez la sortie pour voir la configuration du cgroup et comment elle limite les ressources du processus.

## Indices

- Utilisez la fonction `os.Mkdir` pour créer le répertoire du cgroup.
- Utilisez la fonction `os.WriteFile` pour configurer `memory.max`, `cpu.max`, et `cgroup.procs`.

## Points Clés

- Comprendre comment configurer les cgroups pour limiter l'utilisation de la mémoire et du CPU.
- Apprendre à utiliser le package `os` en Go pour manipuler les cgroups.
- Observer l'effet de la configuration des cgroups sur l'utilisation des ressources du processus.

## Ressources Supplémentaires

- [man cgroups](https://man7.org/linux/man-pages/man7/cgroups.7.html)
- [Documentation du package Go os](https://pkg.go.dev/os)
