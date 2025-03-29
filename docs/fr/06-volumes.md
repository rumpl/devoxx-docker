# Implémentation des Montages de Volumes pour Containers

## Objectif

Apprendre à implémenter la fonctionnalité de montage de volumes pour les containers en utilisant les bind mounts. Cet exercice démontre comment partager des répertoires entre l'hôte et le container, permettant la persistance et le partage des données.

## Étapes

### Étape 1 : Créer la Structure de Répertoires du Volume

1. Configurer les répertoires de volume :
   ```go
   func setupVolume(volumePath, containerPath string) error {
       // TODO:
       // 1. Créer le répertoire source du volume sur l'hôte s'il n'existe pas
       // 2. Créer le point de montage cible dans le container
       // 3. Assurer les permissions appropriées (0755)
       return nil
   }
   ```

### Étape 2 : Implémenter le Bind Mount

1. Créer une fonction pour gérer le bind mounting :

   ```go
   func mountVolume(source, target string) error {
       // Créer le répertoire cible
       if err := os.MkdirAll(target, 0755); err != nil {
           return fmt.Errorf("mkdir %w", err)
       }

       // Effectuer le bind mount
       if err := syscall.Mount(source, target, "", syscall.MS_BIND, ""); err != nil {
           return fmt.Errorf("bind mount %w", err)
       }

       return nil
   }
   ```

### Étape 3 : Ajouter le Démontage de Volume

1. Implémenter le démontage propre des volumes :
   ```go
   func unmountVolume(target string) error {
       // TODO:
       // 1. Démonter le volume en utilisant syscall.Unmount
       // 2. Gérer les erreurs de montage occupé
       // 3. Nettoyer le répertoire du point de montage
       return nil
   }
   ```

### Étape 4 : Intégration avec le Runtime de Container

1. Ajouter la gestion des volumes à votre flux de création de container :

   ```go
   func setupContainerVolumes(containerID string) error {
       volumes := []struct {
           source string
           target string
       }{
           {"/host/path", "/container/path"},
           // Ajouter d'autres mappages de volumes selon les besoins
       }

       for _, vol := range volumes {
           if err := mountVolume(vol.source, vol.target); err != nil {
               return fmt.Errorf("mount volume %s: %w", vol.source, err)
           }
       }

       return nil
   }
   ```

### Étape 5 : Test

1. Tester votre implémentation de volume :

   ```bash
   # Créer des fichiers de test dans le volume de l'hôte
   echo "test data" > /path/to/host/volume/test.txt

   # Exécuter le container avec le volume
   sudo ./container run -v /path/to/host/volume:/container/volume ubuntu /bin/bash

   # Vérifier depuis l'intérieur du container
   cat /container/volume/test.txt
   touch /container/volume/newfile.txt

   # Vérifier que les changements sont visibles sur l'hôte
   ls -l /path/to/host/volume/newfile.txt
   ```

## Indices

- Utilisez `syscall.Mount()` avec le flag `MS_BIND` pour les bind mounts
- Créez toujours les répertoires cibles avant de monter
- N'oubliez pas de gérer le démontage lors du nettoyage du container
- Utilisez `defer` pour les opérations de nettoyage
- Vérifiez les montages existants avant de monter
- Assurez une gestion d'erreurs appropriée et un nettoyage en cas d'échecs

## Points Clés

- Les bind mounts créent une vue d'un répertoire de l'hôte dans le container
- Un nettoyage approprié est essentiel pour éviter les montages orphelins
- Les chemins des volumes doivent exister avant le montage
- Les changements dans les volumes montés sont immédiatement visibles à la fois dans l'hôte et le container
- Les flags de montage affectent le comportement du volume monté

## Ressources Supplémentaires

- [man mount](https://man7.org/linux/man-pages/man2/mount.2.html)
- [man umount](https://man7.org/linux/man-pages/man2/umount.2.html)
- [Linux bind mounts](https://man7.org/linux/man-pages/man8/mount.8.html#BIND_MOUNT_OPERATION)
- [Container volumes](https://docs.docker.com/storage/volumes/)

## Référence des Commandes

### Opérations de Montage

```go
// Bind mount de base
syscall.Mount(source, target, "", syscall.MS_BIND, "")

// Bind mount avec des flags supplémentaires
syscall.Mount(source, target, "", syscall.MS_BIND|syscall.MS_REC, "")

// Démontage
syscall.Unmount(target, 0)
```

### Opérations de Répertoire

```go
// Créer un point de montage
os.MkdirAll(path, 0755)

// Vérifier si le répertoire existe
if _, err := os.Stat(path); os.IsNotExist(err) {
    // Le répertoire n'existe pas
}

// Supprimer le point de montage
os.RemoveAll(path)
```

### Commandes de Débogage

```bash
# Lister les montages
mount | grep container-path

# Vérifier les points de montage
findmnt

# Déboguer les problèmes de montage
dmesg | tail

# Vérifier le namespace de montage
ls -l /proc/$PID/ns/mnt
```

## Exemples de Gestion d'Erreurs

```go
// Gérer un point de montage occupé
if err := syscall.Unmount(target, 0); err != nil {
    if err == syscall.EBUSY {
        // Gérer un point de montage occupé
        return fmt.Errorf("mount point is busy: %w", err)
    }
    return fmt.Errorf("unmount failed: %w", err)
}

// Gérer une source inexistante
if _, err := os.Stat(source); os.IsNotExist(err) {
    return fmt.Errorf("source path does not exist: %w", err)
}
```
