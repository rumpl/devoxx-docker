# Ajout du Support Réseau aux Containers

## Objectif

Dans cet exercice, vous allez implémenter l'isolation réseau et la connectivité pour les containers en utilisant des paires d'ethernet virtuelles (veth) et des namespaces réseau. Cela permettra aux containers de communiquer avec l'hôte et d'accéder à Internet tout en maintenant l'isolation réseau.

## Étapes

### Étape 1 : Créer le Namespace Réseau

1. Ajoutez l'isolation de namespace réseau dans votre code principal de création de container :
   ```go
   cmd.SysProcAttr = &syscall.SysProcAttr{
       // Ajoutez CLONE_NEWNET à vos flags clone existants
       // Cela crée un nouveau namespace réseau pour le container
   }
   ```

### Étape 2 : Implémenter la Création de Paire veth

1. Créez une fonction pour configurer la paire veth :

   ```go
   func SetupVeth(vethName string, pid int) error {
       // TODO: Utilisez les commandes "ip link" pour :
       // 1. Créer une paire veth (veth0 et veth1)
       // 2. Déplacer veth1 vers le namespace réseau du container
       // 3. Configurer veth0 dans le namespace de l'hôte
       // 4. Configurer les règles NAT en utilisant iptables
   }
   ```

2. Créez une fonction de nettoyage pour supprimer la configuration réseau :
   ```go
   func CleanupVeth(vethName string) error {
       // TODO: Nettoyage :
       // 1. Supprimer les règles NAT
       // 2. Supprimer la paire veth
   }
   ```

### Étape 3 : Configurer le Réseau du Container

1. Créez une fonction pour configurer le réseau à l'intérieur du container :
   ```go
   func SetupContainerNetworking(peerName string) error {
       // TODO: À l'intérieur du container :
       // 1. Attribuer une adresse IP à l'interface du container
       // 2. Activer l'interface
       // 3. Configurer la route par défaut
       // 4. Configurer l'interface de loopback
   }
   ```

### Étape 4 : Intégration

1. Ajoutez la configuration réseau à votre flux de création de container :
   ```go
   // Après avoir démarré le processus du container :
   vethName := "veth0"
   if err := SetupVeth(vethName, cmd.Process.Pid); err != nil {
       return err
   }
   defer CleanupVeth(vethName)  // Assurer le nettoyage à la sortie
   ```

### Étape 5 : Test

1. Testez votre implémentation réseau :
   ```bash
   # Depuis l'intérieur du container
   ping 10.0.0.1     # Devrait atteindre l'hôte
   ping 8.8.8.8      # Devrait atteindre Internet
   ping google.com   # Devrait résoudre et atteindre
   ```

## Indices

- Utilisez `exec.Command()` pour exécuter les commandes de configuration réseau
- L'interface du container devrait être dans le sous-réseau 10.0.0.0/24
- Attributions IP courantes :
  - Interface hôte (veth0) : 10.0.0.1
  - Interface container (veth1) : 10.0.0.2
- Les règles iptables requises devraient activer le NAT pour le sous-réseau du container
- N'oubliez pas d'activer le transfert IP sur l'hôte
- Utilisez `defer` pour les opérations de nettoyage

## Points Clés

- Les namespaces réseau fournissent l'isolation réseau
- Les paires veth créent une connexion réseau virtuelle
- Le NAT permet l'accès à Internet depuis le container
- Un nettoyage approprié est essentiel pour éviter les fuites de ressources

## Ressources Supplémentaires

- [man ip-netns](https://man7.org/linux/man-pages/man8/ip-netns.8.html)
- [man veth](https://man7.org/linux/man-pages/man4/veth.4.html)
- [man iptables](https://man7.org/linux/man-pages/man8/iptables.8.html)
- [Linux Network Namespaces](https://man7.org/linux/man-pages/man7/network_namespaces.7.html)
- [Container Networking](https://docs.docker.com/network/)

## Référence des Commandes Requises

### Configuration Réseau de l'Hôte

```bash
# Créer une paire veth
ip link add <veth-host> type veth peer name <veth-container>

# Déplacer l'extrémité du container vers le namespace du container
ip link set <veth-container> netns <pid>

# Configurer l'extrémité de l'hôte
ip addr add 10.0.0.1/24 dev <veth-host>
ip link set <veth-host> up

# Configurer le NAT
iptables -t nat -A POSTROUTING -s 10.0.0.0/24 -j MASQUERADE
```

### Configuration Réseau du Container

```bash
# Configurer l'interface du container
ip addr add 10.0.0.2/24 dev <veth-container>
ip link set <veth-container> up
ip link set lo up
ip route add default via 10.0.0.1
```
