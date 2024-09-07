.PHONY: husky format

husky:
	chmod +x formatter.sh; chmod +x pre-commit.sh; cp pre-commit.sh .git/hooks/pre-commit

format:
	./formatter.sh .
