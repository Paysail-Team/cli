#! /usr/bin/env node

import figlet, { fonts } from 'figlet'
import chalk from 'chalk';
import * as commander from 'commander';
import axios from 'axios';
import fs from 'fs';
import path from 'path';

const program = new commander.Command();
program
    .version("1.0.0")
    .description(`${logo()}\n${chalk.whiteBright(`Syro CLI allows you to access your secrets and inject them into your CI/CD pipelines.`)}`)
    
program
    .command('pull')
    .description('Pulls all secrets from the given project and creates a .env file in the current directory')
    .argument('<accessToken>', 'The access token')
    .argument('[environment]', 'The target environment. Defaults to production', "production")
    .action(pull);

program.parse();

if (!process.argv.slice(2).length) {
    program.outputHelp();
}

function logo() {
    return `
${chalk.bgWhiteBright(chalk.black(`
  ┏┓┓┏┳┓┏┓  
  ┗┓┗┫┣┫┃┃  
  ┗┛┗┛┛┗┗┛  
`))}`
}

function host(environment: string) {
    if (environment && environment.length > 0) {
        switch (environment) {
            case 'production': return "https://api-production.syro.com"
            case 'staging': return "https://api-staging.syro.com"
            case 'local': return "http://localhost:1400"
            default: return "https://api-production.syro.com"
        }
    } else {
        return "https://api-production.syro.com"
    }
}

function maskedValue(value: string) {
    const defaultMask = '•••••'
    if (value.length === 0) {
        return defaultMask
    }

    function generate(value: string, numberOfVisibleCharacters: number): string {
        if (value.length === 1) {
            return defaultMask
        }

        if (numberOfVisibleCharacters <= 0) {
            return defaultMask
        }

        if (value.length <= numberOfVisibleCharacters) {
            return generate(value, numberOfVisibleCharacters - 1)
        }
        return `${defaultMask}${value.substring(value.length - numberOfVisibleCharacters, value.length)}`
    }

    return generate(value, 5)
}

async function pull(accessToken: string, environment: string) {
    try {
        console.log(logo())
        console.log(chalk.yellowBright(`Pulling secrets...`))
        const response = await axios({
            "method": "POST",
            "url": `${host(environment)}/secrets`,
            "headers": {
                "Content-Type": "application/json; charset=utf-8"
            },
            "data": {
                "accessToken": accessToken
            }
        })
        const secrets: { key: string, value: string }[] = response.data.result
        if (secrets === undefined) {
            throw new Error("0")
        }
        const secretsTable = secrets.map(secret => {
            return { key: secret.key, value: maskedValue(secret.value) };
        })
        console.log(chalk.green(`✔ Pulled ${secrets.length} secret${secrets.length === 1 ? "" : "s"}.`))
        console.table(secretsTable)

        console.log(chalk.yellowBright(`Generating .env file...`))
        let fileText = ""
        secrets.forEach(item => {
            fileText = fileText + `${item.key}='${item.value}'\n`
        })
        fs.writeFileSync(`.env`, fileText)
        console.log(chalk.green(`✔ Generated .env at ${path.resolve(process.cwd())}/.env\n`))

    } catch (error: any) {
        console.log(chalk.redBright(`
⨯ Unable to pull secrets. Check accessToken and try again.\n\nError: ${error}\n`))
    }
}