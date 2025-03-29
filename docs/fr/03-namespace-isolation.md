# Implémentation de l'Isolation par Namespaces

## Objectif

Dans cet exercice, vous apprendrez à implémenter l'isolation par namespaces et la configuration du nom d'hôte pour les containers. Vous vous concentrerez sur la création d'un nouveau namespace PID et la définition d'un nom d'hôte personnalisé à l'aide du namespace UTS.

## Étapes

### Étape 1 : Ajouter l'Isolation par Namespaces

1.  Modifier la création du processus parent pour inclure les flags de namespaces :
    ```go
    func run() error {
    cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args...)...)

             //TODO:
            // 1. Ajouter les flags de namespaces pour les namespaces PID et UTS

            if err := cmd.Wait(); err != nil {
               return fmt.Errorf("wait %w", err)
            }

            fmt.Printf("Container exited with exit code %d\n", cmd.ProcessState.ExitCode())
        }
        ```

    <details>
    <summary>Indice</summary>
    examinez la structure `syscall.SysProcAttr`
    </details>

### Étape 2 : Implémenter les Changements de Nom d'Hôte

1. Ajouter la configuration du nom d'hôte au processus enfant :
`go
    func child() error {
        //TODO: 
		// 1. Définir le nom d'hôte du container
		// 2. Afficher le PID de l'enfant 
		// 3. Afficher le nom d'hôte pour vérifier le changement
    }
    `
<details>
<summary>Indice</summary>
examinez la fonction `syscall.Sethostname`
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

Nous avons maintenant implémenté l'isolation des namespaces PID et UTS, fournissant une isolation des processus et une configuration personnalisée du nom d'hôte pour les containers.  
C'est une étape cruciale vers la construction d'un runtime de container pleinement fonctionnel.

[Étape suivante](04-namespaces-and-chroot.md)

## Indices

- Utilisez `syscall.CLONE_NEWPID` pour l'isolation du namespace PID
- Utilisez `syscall.CLONE_NEWUTS` pour l'isolation du namespace de nom d'hôte
- Des privilèges root sont nécessaires pour les opérations sur les namespaces
- Le processus enfant doit avoir le PID 1 dans son namespace

## Points Clés

- Le namespace PID fournit l'isolation des processus
- Le namespace UTS permet un nom d'hôte personnalisé
- Les changements de namespace nécessitent des privilèges root
- Le processus enfant se voit comme ayant le PID 1

## Ressources Supplémentaires

- [man namespaces](https://man7.org/linux/man-pages/man7/namespaces.7.html)
- [man clone](https://man7.org/linux/man-pages/man2/clone.2.html)
- [Package Go syscall](https://pkg.go.dev/syscall)

## Référence des Commandes

### Opérations sur les Namespaces

```go
// Créer de nouveaux namespaces
cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
}

// Définir le nom d'hôte
syscall.Sethostname([]byte("new-hostname"))
```

### Commandes de Débogage

```bash
# Vérifier les namespaces du processus
ls -l /proc/$$/ns/

# Voir le nom d'hôte
hostname

# Vérifier le PID dans différents namespaces
ps aux
```

### Exemples de Gestion d'Erreurs

```go
// Gérer les erreurs de nom d'hôte
if err := syscall.Sethostname([]byte("container-host")); err != nil {
    if os.IsPermission(err) {
        return fmt.Errorf("permission denied: run with sudo: %w", err)
    }
    return fmt.Errorf("hostname error: %w", err)
}
```
