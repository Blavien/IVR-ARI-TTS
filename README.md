# IVR-ARI-TTS

Aplicação que utiliza a ARI API para comunicar com o servidor Asterisk.
Baseada no exemplo Play do repositório https://github.com/CyCoreSystems/ari/tree/master/_examples
É utilizado a lib da Google para fazer o request de TTS e gerar uma sintese de aúdio.

Não está a funcionar o envio do audio gerado para o servidor.
Iremos fazer uso da feature do Asterisk chamada ExternalMedia para que possamos fazer stream do audio para o canal de comunicação.

ExternalMedia - https://github.com/CyCoreSystems/ari/blob/master/channel.go
