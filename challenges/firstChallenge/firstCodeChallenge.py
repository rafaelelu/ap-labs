test = input('Input: ')

first = ''
i = 0
while (test[i] not in first and i < len(test)):
    first  += test[i]
    i += 1
second = ''
while (i < len(test)):
    if  (test[i] in second):
        second  = test[i]
    else:
        second += test[i]
    if len(second) > len(first):
        first = second
        second = ''
    i += 1

print('Output: ',len(first))
