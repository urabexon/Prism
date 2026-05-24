# Prism 💍
> Une TUI pour réviser les Pull Requests GitHub dans votre terminal

## Fonctionnalités
- Révisez les PRs sans quitter le terminal
- Raccourcis clavier de style Vim pour une navigation rapide
- Affichage des diffs, commentaires et fusion
- Vérification du statut CI
- Sauvegarde automatique de la progression

## Installation
```bash
go install github.com/urabexon/prism@latest
```

## Utilisation
1. Lancer dans le dépôt actuel
```bash
prism
```
2. Spécifier un dépôt
```bash
prism owner/repo
```

## Raccourcis clavier
### Liste des PRs
| Touche  | Action          |
| ------- | --------------- |
| `j/k`   | Naviguer        |
| `Enter` | Ouvrir la PR    |
| `m`     | Fusionner       |
| `c`     | Vérifications CI|
| `q`     | Quitter         |

### Vue des diffs
| Touche  | Action              |
| ------- | ------------------- |
| `j/k`   | Défiler             |
| `n/N`   | Hunk suivant/préc.  |
| `V`     | Sélection visuelle  |
| `space` | Marquer comme révisé|
| `esc`   | Retour              |

## Prérequis
- Go 1.21+
- GitHub CLI (gh) authentifié

## Licence
MIT
