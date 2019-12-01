*Algoritmo Dijkstra-Scholten*:
	- *Número de processos*: 4

	- *Comunicação*: RPC, copiado dos outros labs

	- *Estados dos processos*: Active, Passive

	- *Parâmetros*:
		- -init, T/F:
			- para quem é true, começa o processamento ativo e espera uns segundos para deixar os outros processos se tornarem online
			- ele quem inicia o processamento básico e a detecção de término

	- *Mensagem básica*:
		- Remetente

	- *Mensagem control*:
		- Remetente

	- *Mensagem finish*:
	 	- Remetente

	- *Ao p receber uma mensagem básica*:
		- se está active:
			- responde com um control		

		- se está passive:
			- fica active
			- guarda o seu pai

		- escolhe um número aleatório de tempo
		- printa "Estarei ativo por x tempo"
		- dorme por x tempo
		- aleatoriamente escolhe se vai continuar a atividade
		- Se sim,
			- Escolhe um número aleatório n de processos para os quais vai mandar uma mensagem básica,
			- Para cada um dos n processos:
				- Manda a mensagem
				- p adiciona q a seus filhos na árvore
			- Fica passive
			- Checa se ele ainda tem filhos
			- se p não tiver mais filhos
				- se p não tem pai, declara terminação
				- se p tem pai, p manda um control para seu pai
		- Se não
			- Fica quiet

	- *Quando p ficar inativo*:
		- espera um determinado timeout antes de tentar detectar término
		- printa "inicia WAVE <p, t>" e faz broadcast de mensagem wave

	- *Quando p recebe um control de q*
		- p remove q de seus filhos na árvore
		- se p não tiver mais filhos e estiver passivo:
			- se p não tem pai, declara terminação
			- se p tem pai, p manda um leave para seu pai
