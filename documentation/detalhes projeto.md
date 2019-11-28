# Algoritmo de Rana:
- *Número de processos*: 4

- *Comunicação*: RPC, copiado dos outros labs

- *Estados dos processos*: Active, Passive, Quiet

- *Parâmetros*:
	- active, T/F:
	- para quem é true, começa o processamento ativo e espera uns segundos para deixar os outros processos se tornarem online

- *Mensagem básica*:
	- Remetente
	- Clock

- *Mensagem ack*:
	- Remetente
	- Clock

- *Mensagem WAVE*:
	- Iniciador
	- Clock
	- Vetor bool de assinaturas

- *Ao p receber uma mensagem básica ou começar ativo*:
	- fica active
	- responde ack(p,t)
	- escolhe um número aleatório de tempo
	- printa "Estarei ativo por x tempo"
	- dorme por x tempo
	- aleatoriamente escolhe se vai continuar a atividade
	- Se sim,
		- Escolhe um número aleatório n de processos para os quais vai mandar uma mensagem básica,
		- Manda as n mensagens
		- Fica passive e espera n acks
		- Ao receber n acks, fica quiet
	- Se não
		- Fica quiet

- *Quando p ficar inativo*:
	- espera um determinado timeout antes de tentar detectar término
	- printa "inicia WAVE <p, t>" e faz broadcast de mensagem wave

- *Quando q recebe mensagem wave(p,t)*:
	- se q está active,
		- printa "rejeita WAVE<p,t>, ainda ativo"
	- se q está passive,
		- printa "rejeita WAVE<p,t>, esperando acks"
	- se q está quiet e t < Tq,
		- printa "rejeita WAVE<p,t>, clock incompatível"
	- se q está quiet e t >= Tq,
		- printa "aceita WAVE<p,t>"
		- copia WAVE<p,t> e faz assinaturas[q] = True
		- faz broadcast da WAVE<p,t>
