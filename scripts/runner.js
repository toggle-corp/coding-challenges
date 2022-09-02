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

function getTestData() {
    const testInputStr = fs.readFileSync(TEST_INPUT, 'utf8');
    const testInpItems = testInputStr.split(/\r?\n/).filter(x => !!x.trim());

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
    const testInput = testInpItems.map(lineStr => eval(enclose(lineStr)));
    const testOutputStr = fs.readFileSync(TEST_OUTPUT, 'utf8');
    const testOutput = testOutputStr.split(/\r?\n/).filter(x => !!x.trim()).map(lineStr => eval(lineStr))
    return [testInput, testOutput, inpApplyFn]
}

function main() {
    const [inps, expecteds, inpApplyFn] = getTestData();
    assert(inps.length == expecteds.length);

    let passedCount = 0;
    inps.forEach((inp, i) => {
        const result = inpApplyFn(solution.solution, inp);
        const expected = expecteds[i];
        if (result != expected) {
            console.error(`Failed for test case ${i+1}.`);
            console.error(`Expected: ${expected} Got: ${result}.`);
            console.error(`Passed ${passedCount}/${inps.length}`);
            process.exit(43);
        }
        passedCount += 1;
    });
}

main();