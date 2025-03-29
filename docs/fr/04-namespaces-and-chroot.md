# Gestion des Namespaces et du Répertoire Racine

## Objectif

Apprendre à gérer l'isolation du système de fichiers dans un environnement conteneurisé en implémentant des namespaces de montage et en changeant le répertoire racine à l'aide de `chroot`.  
Cet exercice démontre comment créer un environnement de système de fichiers isolé.

## Étapes

### Étape 1 : Ajouter le Namespace de Montage

1. Modifier le processus parent pour inclure la capacité de namespace de montage au processus enfant :
   ```go
   func parent() error {
       cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[1:]...)...)

       // TODO:
       // 1. Ajouter le flag de namespace de montage
   }
   ```

### Étape 2 : Configurer la Structure du Répertoire Racine

1.  Créer la fonction pour configurer le système de fichiers racine :
    ```go
    func child() error {
    // TODO:
    // 1. Créer le répertoire racine de base à "/fs/container/rootfs"
    // 2. Définir les permissions appropriées (0755)
    // 3. Gérer toutes les erreurs potentielles

            return nil
        }
        ```

    <details>
    <summary>Indice</summary>
    examinez la fonction `os.MkdirAll`
    </details>

### Étape 3 : Changer le Répertoire Racine

1.  Implémenter la configuration du répertoire racine du container :
    ```go
    func setupContainer() error {
    // TODO:
    // 1. Afficher le répertoire de travail actuel
    // 2. Changer la racine vers "/fs/container/rootfs"
    // 3. Changer le répertoire courant vers la racine ("/")
    // 4. Gérer toutes les erreurs potentielles
    // 5. Implémenter une gestion d'erreurs appropriée
    // 6. Afficher le nouveau répertoire de travail

            return nil
        }
        ```

    <details>
    <summary>Indice</summary>
    examinez les fonctions `syscall.Chroot` et `os.Chdir`

</details>

### Étape 4 : Test

1. Compiler et exécuter votre programme :

   ```bash
   # Compiler le programme
   make

   # Exécuter avec sudo (nécessaire pour les opérations sur les namespaces)
   sudo ./bin/devoxx-container
   ```

### Résumé

Nous avons maintenant implémenté l'isolation du namespace de montage et changé le répertoire racine pour le container.  
Cela fournit un environnement de système de fichiers isolé pour le container.

[Étape suivante](05-cgroups.md)

## Indices

- Utilisez `syscall.CLONE_NEWNS` pour l'isolation du namespace de montage
- Des privilèges root sont nécessaires pour les opérations sur les namespaces
- Utilisez des chemins absolus lors de la manipulation des répertoires
- N'oubliez pas de gérer le nettoyage en cas d'erreurs
- Vérifiez si les répertoires existent avant les opérations
- Utilisez `defer` pour les opérations de nettoyage

## Points Clés

- Les namespaces de montage fournissent l'isolation du système de fichiers
- `chroot` change la vue du répertoire racine
- Un nettoyage approprié est essentiel pour éviter les fuites de ressources
- Les opérations sur les namespaces nécessitent une gestion d'erreurs minutieuse

## Ressources Supplémentaires

- [man mount_namespaces](https://man7.org/linux/man-pages/man7/mount_namespaces.7.html)
- [man chroot](https://man7.org/linux/man-pages/man2/chroot.2.html)
- [Linux Filesystem Hierarchy Standard](https://refspecs.linuxfoundation.org/FHS_3.0/fhs/index.html)

## Référence des Commandes

### Opérations sur les Namespaces

```go
// Créer un namespace de montage
syscall.CLONE_NEWNS

// Changer la racine
syscall.Chroot(path)
```

### Commandes de Débogage

```bash
# Vérifier la structure du système de fichiers
ls -la /fs/container/rootfs

# Voir les namespaces de montage
ls -l /proc/$$/ns/mnt

# Voir les namespaces du processus
ls -l /proc/$$/ns/
```

### Exemples de Gestion d'Erreurs

```go
// Gérer les erreurs de chroot
if err := syscall.Chroot("/fs/container/rootfs"); err != nil {
    if os.IsPermission(err) {
        return fmt.Errorf("chroot permission denied (run with sudo): %w", err)
    }
    return fmt.Errorf("chroot failed: %w", err)
}

// Gérer les opérations de répertoire
if err := os.MkdirAll("/path/to/dir", 0755); err != nil {
    return fmt.Errorf("failed to create directory: %w", err)
}
```
