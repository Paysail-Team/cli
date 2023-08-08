#! /usr/bin/env node
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
import chalk from 'chalk';
import * as commander from 'commander';
import axios from 'axios';
import fs from 'fs';
import path from 'path';
const program = new commander.Command();
program
    .version("1.0.0")
    .description(`${logo()}\n${chalk.whiteBright(`Syro CLI allows you to access your secrets and inject them into your CI/CD pipelines.`)}`)
    .command('pull')
    .description('Pulls all secrets from the given project and creates a .env file in the current directory')
    .argument('<accessToken>', 'The access token')
    .argument('<projectId>', 'The project id')
    .option('-e, --env', 'The target environment', "production")
    .action(pull);
program.parse();
if (!process.argv.slice(2).length) {
    program.outputHelp();
}
var Environment;
(function (Environment) {
    Environment["production"] = "production";
    Environment["staging"] = "staging";
    Environment["local"] = "local";
})(Environment || (Environment = {}));
function logo() {
    return `
${chalk.bgWhiteBright(chalk.black(`
  ┏┓┓┏┳┓┏┓  
  ┗┓┗┫┣┫┃┃  
  ┗┛┗┛┛┗┗┛  
`))}`;
}
function host() {
    const options = program.opts();
    if (options.env) {
        switch (options.env) {
            case Environment.production: return "https://api-production.syro.com";
            case Environment.staging: return "https://api-production.syro.com";
            case Environment.local: return "http://localhost:1400/";
            default: return "https://api-production.syro.com";
        }
    }
    else {
        return "https://api-production.syro.com";
    }
}
function maskedValue(value) {
    const defaultMask = '•••••';
    if (value.length === 0) {
        return defaultMask;
    }
    function generate(value, numberOfVisibleCharacters) {
        if (value.length === 1) {
            return defaultMask;
        }
        if (numberOfVisibleCharacters <= 0) {
            return defaultMask;
        }
        if (value.length <= numberOfVisibleCharacters) {
            return generate(value, numberOfVisibleCharacters - 1);
        }
        return `${defaultMask}${value.substring(value.length - numberOfVisibleCharacters, value.length)}`;
    }
    return generate(value, 5);
}
function pull(accessToken, projectId) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            console.log(logo());
            console.log(chalk.yellowBright(`Pulling secrets...`));
            const response = yield axios({
                "method": "POST",
                "url": `${host()}/secrets`,
                "headers": {
                    "Content-Type": "application/json; charset=utf-8"
                },
                "data": {
                    "accessToken": accessToken,
                    "projectId": projectId
                }
            });
            const secrets = response.data.result;
            if (secrets === undefined) {
                throw new Error("0");
            }
            const secretsTable = secrets.map(secret => {
                return { key: secret.key, value: maskedValue(secret.value) };
            });
            console.log(chalk.green(`✔ Pulled ${secrets.length} secret${secrets.length === 1 ? "" : "s"}.`));
            console.table(secretsTable);
            console.log(chalk.yellowBright(`Generating .env file...`));
            let fileText = "";
            secrets.forEach(item => {
                fileText = fileText + `${item.key}='${item.value}'\n`;
            });
            fs.writeFileSync(`.env`, fileText);
            console.log(chalk.green(`✔ Generated .env at ${path.resolve(process.cwd())}/.env\n`));
        }
        catch (error) {
            console.log(chalk.redBright(`
⨯ Unable to pull secrets. Check accessToken and projectId, and try again.\n\nError: ${error}\n`));
        }
    });
}
//# sourceMappingURL=index.js.map