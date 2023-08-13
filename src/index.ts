#! /usr/bin/env node

import chalk from 'chalk';
import * as commander from 'commander';
import axios from 'axios';
import fs from 'fs';
import path from 'path';
import confirm from '@inquirer/confirm';

const program = new commander.Command();
program
    .version("1.0.0")
    .description(`${logo()}\n${chalk.whiteBright(`Syro CLI allows you to access your secrets and inject them into your CI/CD pipelines.`)}`)

program.command('pull')
    .description('Pulls all secrets from the given project and creates a .env file in the current directory')
    .argument('<accessToken...>', 'The access token')
    .option('-e, --env [env]', 'The target environment', 'production')
    .option('-n, --filename [filename]', 'The target environment file name')
    .option('-f, --force', 'Create the environment file without protections')
    .action((accessTokens: string[], options: { [index: string]: string }) => {
        return pull(accessTokens, options)
    });

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

async function pull(accessTokens: string[], options: { [index: string]: string }) {
    try {
        const environment = options['env']
        const isForced = options['force']

        if (environment === undefined) {
            throw new Error("1000")
        }

        console.log(logo())
        console.log(chalk.yellowBright(`Pulling secrets...`))

        if (accessTokens === undefined) {
            throw new Error("1001")
        }

        if (accessTokens.length === 0) {
            throw new Error("1002")
        }

        const response = await axios({
            "method": "POST",
            "url": `${host(environment)}/cli/secrets`,
            "headers": {
                "Content-Type": "application/json; charset=utf-8"
            },
            "data": {
                "accessTokens": accessTokens
            }
        })
        const data: { pn: string[], efn: string[], i: { key: string, value: string }[], cefn: string } = response.data

        if (data === undefined || data.i === undefined || data.pn === undefined || data.efn === undefined || data.cefn === undefined) {
            throw new Error("1003")
        }

        const targetFileName: string = options['filename'] ?? data.cefn

        const secretsTable = data.i.map(secret => {
            return { key: secret.key, value: maskedValue(secret.value) };
        })

        const display_projectNames = data.pn.length === 1 ? data.pn[0] : data.pn.length > 4 ? `${data.pn.length} projects` : `${data.pn.slice(0, data.pn.length - 1).slice(0, 4).join(", ")} and ${data.pn[data.pn.length - 1]}`

        console.log(chalk.green(`✔ Pulled ${data.i.length} secret${data.i.length === 1 ? "" : "s"} from ${display_projectNames}.`))
        console.table(secretsTable)

        console.log(chalk.yellowBright(`Generating ${targetFileName} file...`))

        let didOverwrite = false
        if (isForced === undefined && fs.existsSync(targetFileName)) {
            const answer = await confirm({ message: chalk.whiteBright(`${targetFileName} already exists in current directory. Do you want to overwrite the existing file?`) })
            if (answer === false) {
                console.log(chalk.green(`✔ Existing ${targetFileName} was not overwritten.\n`))
                return
            }
            didOverwrite = answer
        }

        let fileText = ""
        data.i.forEach(item => {
            fileText = fileText + `${item.key}='${item.value}'\n`
        })
        fs.writeFileSync(targetFileName, fileText)

        console.log(chalk.green(`✔ ${didOverwrite ? "Overwrote" : "Generated"} ${targetFileName} at ${path.resolve(process.cwd())}/\n`))
    } catch (error: any) {
        const display_accessToken = `accessToken${accessTokens.length === 1 ? '' : 's'}`
        console.log(chalk.redBright(`\n⨯ Unable to pull secrets. Please check your ${display_accessToken} and try again.`))
        console.log(chalk.red(`  ${error}\n`))
    }
}