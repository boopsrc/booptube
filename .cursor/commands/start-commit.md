# start-commit

Analise tudo que foi feito e defina um nome para o commit seguindo o padrao

A mensagem de commit deve ser estruturada da seguinte forma:

<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
O commit contém os seguintes elementos estruturais, para comunicar a intenção aos usuários da sua biblioteca:

correção: um commit desse tipo fix corrige um bug em sua base de código (isso se correlaciona com PATCHa versão semântica).
recurso: um commit desse tipo feat introduz uma nova funcionalidade à base de código (isso se correlaciona com MINORa versão semântica).
ALTERAÇÃO QUE QUEBRA A COMPATIBILIDADE: um commit que possui um rodapé BREAKING CHANGE:ou adiciona um !após o tipo/escopo introduz uma alteração de API que quebra a compatibilidade (correspondente a MAJORno Versionamento Semântico). Uma ALTERAÇÃO QUE QUEBRA A COMPATIBILIDADE pode fazer parte de commits de qualquer tipo .
Tipos diferentes de `std::string` e fix:` std :: vector` feat:são permitidos; por exemplo , @commitlint/config-conventional (baseado na convenção do Angular ) recomenda build:` chore:std ::string` , `std::vector` , `std :: string ...ci:docs:style:refactor:perf:test:
Rodapés diferentes BREAKING CHANGE: <description>podem ser fornecidos e seguem uma convenção semelhante ao formato de trailer do Git .
Tipos adicionais não são obrigatórios pela especificação de Commits Convencionais e não têm efeito implícito no Versionamento Semântico (a menos que incluam uma ALTERAÇÃO QUE QUEBRA A COMPATIBILIDADE). Um escopo pode ser fornecido ao tipo de um commit para fornecer informações contextuais adicionais e é contido entre parênteses, por exemplo, feat(parser): add ability to parse arrays.

Exemplos
Mensagem de commit com descrição e rodapé de alteração incompatível
feat: allow provided config object to extend other configs

BREAKING CHANGE: `extends` key in config file is now used for extending other config files
Mensagem de compromisso !para chamar a atenção para a mudança radical
feat!: send an email to the customer when a product is shipped
A mensagem de compromisso deve incluir o escopo e !chamar a atenção para mudanças significativas.
feat(api)!: send an email to the customer when a product is shipped
Mensagem de commit com !rodapé "ALTERAÇÃO QUEBRA A QUALIDADE"
feat!: drop support for Node 6

BREAKING CHANGE: use JavaScript features not available in Node 6.
Mensagem de commit sem corpo
docs: correct spelling of CHANGELOG
Mensagem de commit com escopo
feat(lang): add Polish language
Mensagem de commit com corpo de vários parágrafos e vários rodapés
fix: prevent racing of requests

Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.

Reviewed-by: Z
Refs: #123


comando a ser executado 
git add .
git commit -m "mensagem aqui"
