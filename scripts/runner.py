# The solution module will be available when the runner script is invoked
from solution import solution

import sys


TEST_INPUT = './test_inputs'
TEST_OUTPUT = './test_outputs'

def main(delimeter = '\n'):
    # NOTE: only used for input
    with open(TEST_INPUT) as finp, open(TEST_OUTPUT) as fout:
        inpstr = finp.read()
        splitted = [x.strip() for x in inpstr.split(delimeter)]
        inputs = [x for x in splitted if x]
        outputs = [x for x in fout.readlines() if x.strip()]
    assert len(inputs) == len(outputs), "Test cases not aligned"

    total_count = len(inputs)
    passed_count = 0
    for i, (inp, out) in enumerate(zip(inputs, outputs)):
        try:
            inpargs = eval(inp)
        except SyntaxError:  # Just treat input as string
            inpargs = inp

        try:
            expected = eval(out)
        except SyntaxError:
            expected = out

        if isinstance(inpargs, tuple):
            # Means solution takes multiple params
            result = solution(*inpargs)
        else:
            result = solution(inpargs)

        if result != expected:
            print(f'Failed for test case {i+1}.', file=sys.stderr)
            print(f'Passed {passed_count}/{total_count}', file=sys.stderr)
            sys.exit(43)

        passed_count += 1  # TODO not used for now


if __name__ == '__main__':
    if len(sys.argv) >= 2:
        delimeter = sys.argv[1]
        main(delimeter)
    else:
        main()
