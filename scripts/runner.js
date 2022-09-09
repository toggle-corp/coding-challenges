const TEST_INPUT = './test_inputs';
const TEST_OUTPUT = './test_outputs';

const solution = require('./solution');
const fs = require('fs');
const assert = require('assert');

function expApply(f, x) {
    return f(...x);
}
function normalApply(f,x) {
    return f(x);
}

function tryEval(x) {
    try {
        return eval(x);
    } catch (e) {
        return x;
    }
}

function getTestData(delimeter='\n') {
    const testInputStr = fs.readFileSync(TEST_INPUT, 'utf8');
    const testInpItems = testInputStr.split(delimeter).filter(x => !!x.trim());

    assert(testInpItems.length > 0);

    const first = testInpItems[0];
    // Check if the input are multiple arguments or single argument
    // TODO: this is just a crude check
    let inpApplyFn;
    let shouldEncloseBrackets = false;
    if (first.trim().includes(',') && first.trim()[0] != '[') {
        inpApplyFn = expApply;
        shouldEncloseBrackets = true;
    } else {
        inpApplyFn = normalApply;
    }
    const enclose = (x) => shouldEncloseBrackets ? '['+x+']' : x;
    const testInput = testInpItems.map(lineStr => tryEval(enclose(lineStr)));
    const testOutputStr = fs.readFileSync(TEST_OUTPUT, 'utf8');
    const testOutput = testOutputStr.split(/\r?\n/).filter(x => !!x.trim()).map(lineStr => tryEval(lineStr))
    return [testInput, testOutput, inpApplyFn]
}

function main(delimeter='\n') {
    const [inps, expecteds, inpApplyFn] = getTestData(delimeter);
    assert(inps.length == expecteds.length);

    let passedCount = 0;
    inps.forEach((inp, i) => {
        const result = inpApplyFn(solution.solution, inp);
        const expected = expecteds[i];
        if (JSON.stringify(result) != JSON.stringify(expected)) {
            console.error(`Failed for test case ${i+1}.`);
            console.error(`Passed ${passedCount}/${inps.length}`);
            process.exit(43);
        }
        passedCount += 1;
    });
}

args = process.argv;
if (args.length > 2) {
    main(args[2]);
} else {
    main();
}
