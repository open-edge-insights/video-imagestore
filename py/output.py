def handleOut(code, value):
    return {
        'NotSupported': (False, 'Currently System Not Supports  : ' + value),
        'BaseImplError': (False, 'Error from Base Implementation : ' + value),
        'success': (True, value),
        'error': (False, value),
    }[code]
