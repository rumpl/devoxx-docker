# Création de Processus Parents et Enfants

## Objectif

Techniquement, un container est un processus Linux isolé à l'aide de `cgroups` (control
groups) et de `namespaces` (comme PID, net, mount, user, etc.) pour restreindre l'utilisation
des ressources et fournir une vue séparée du système.

Commençons petit et créons d'abord un nouveau processus.

## Étapes

### Étape 1 : Créer la Fonction Principale

1. Établir la structure de base du programme :
   ```go
   func main() {
      //TODO: Vérifier si nous exécutons la commande initiale ou le processus enfant
      // Si les arguments contiennent "child", appeler child()
      // Sinon, continuer avec la création du processus parent
   }
   ```

### Étape 2 : Implémenter le Processus Enfant

1. Créer le gestionnaire de processus enfant :
   ```go
   func child() error {
      //TODO:
      // 1. Afficher le PID actuel pour démontrer l'isolation du namespace
      // 2. Exécuter la commande souhaitée, un simple affichage `Hello from child` dans la console est suffisant pour l'instant
      // 3. Maintenir le processus en cours d'exécution pour observer l'isolation
   }
   ```

### Étape 3 : Implémenter la Création du Processus Parent

1. Créer une fonction pour gérer la logique du processus parent :
`go
    func run() error {
       //TODO: 
       // 1. Créer une nouvelle commande en utilisant l'exécutable actuel
       // 2. Configurer stdin/stdout/stderr
       // 3. Démarrer le processus enfant
       // 4. Attendre la fin et afficher un message nous informant que le processus enfant s'est terminé
    }
    `
<details>
<summary>Indices</summary>

Utilisez `/proc/self/exe` pour ré-exécuter le même processus

</details>

### Étape 4 : Test

1. Compiler et exécuter votre programme :

   ```bash
   # Compiler le programme
   go build -o devoxx-container

   # Exécuter le programme
   ./devoxx-container
   ```

### Résumé

Nous avons la première étape de base dans notre voyage pour créer un container, nous avons un
processus parent qui peut gérer le processus enfant. Ce processus enfant deviendra bientôt
un véritable container.

[Étape suivante](03-namespace-isolation.md)

## Indices

- Utilisez `/proc/self/exe` pour obtenir le chemin vers l'exécutable actuel
- Utilisez `os.Args` pour détecter si le programme s'exécute en tant qu'enfant
- N'oubliez pas de gérer toutes les erreurs potentielles
- Utilisez `cmd.Start()` et `cmd.Wait()` pour un meilleur contrôle des processus

## Points Clés

- Le processus parent crée et gère le processus enfant
- Le processus enfant s'exécute avec son propre PID
- Une gestion appropriée des erreurs est cruciale pour la gestion des processus
- Les flux d'E/S standard doivent être correctement connectés

## Ressources Supplémentaires

- [Package Go os/exec](https://pkg.go.dev/os/exec)
- [Package Go os](https://pkg.go.dev/os)

## Référence des Commandes

### Informations sur les Processus

```bash
# Voir l'arborescence des processus
ps -ef --forest

# Obtenir les informations sur le processus actuel
ps -p $$

# Voir l'environnement du processus
ps eww -p <pid>
```

### Opérations Courantes

```go
// Obtenir le PID actuel
pid := os.Getpid()

// Obtenir le PID parent
ppid := os.Getppid()

// Créer une commande avec des arguments
cmd := exec.Command("program", "arg1", "arg2")

// Exécuter la commande et attendre la fin
err := cmd.Run()
```
