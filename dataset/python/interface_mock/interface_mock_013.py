import re
from pygments.lexer import Lexer, RegexLexer, do_insertions, bygroups, include, default, this, using, words, line_re
from pygments.token import Punctuation, Whitespace, Text, Comment, Operator, Keyword, Name, String, Number, Generic
from pygments.util import shebang_matches
__all__ = ['BashLexer', 'BashSessionLexer', 'TcshLexer', 'BatchLexer', 'SlurmBashLexer', 'MSDOSSessionLexer', 'PowerShellLexer', 'PowerShellSessionLexer', 'TcshSessionLexer', 'FishShellLexer', 'ExeclineLexer']

class BashLexer(RegexLexer):
    name = 'Bash'
    aliases = ['bash', 'sh', 'ksh', 'zsh', 'shell', 'openrc']
    filenames = ['*.sh', '*.ksh', '*.bash', '*.ebuild', '*.eclass', '*.exheres-0', '*.exlib', '*.zsh', '.bashrc', 'bashrc', '.bash_*', 'bash_*', 'zshrc', '.zshrc', '.kshrc', 'kshrc', 'PKGBUILD']
    mimetypes = ['application/x-sh', 'application/x-shellscript', 'text/x-shellscript']
    url = 'https://en.wikipedia.org/wiki/Unix_shell'
    version_added = '0.6'
    tokens = {'root': [include('basic'), ('`', String.Backtick, 'backticks'), include('data'), include('interp')], 'interp': [('\\$\\(\\(', Keyword, 'math'), ('\\$\\(', Keyword, 'paren'), ('\\$\\{#?', String.Interpol, 'curly'), ('\\$[a-zA-Z_]\\w*', Name.Variable), ('\\$(?:\\d+|[#$?!_*@-])', Name.Variable), ('\\$', Text)], 'basic': [('\\b(if|fi|else|while|in|do|done|for|then|return|function|case|select|break|continue|until|esac|elif)(\\s*)\\b', bygroups(Keyword, Whitespace)), ('\\b(alias|bg|bind|builtin|caller|cd|command|compgen|complete|declare|dirs|disown|echo|enable|eval|exec|exit|export|false|fc|fg|getopts|hash|help|history|jobs|kill|let|local|logout|popd|printf|pushd|pwd|read|readonly|set|shift|shopt|source|suspend|test|time|times|trap|true|type|typeset|ulimit|umask|unalias|unset|wait)(?=[\\s)`])', Name.Builtin), ('\\A#!.+\\n', Comment.Hashbang), ('#.*\\n', Comment.Single), ('\\\\[\\w\\W]', String.Escape), ('(\\b\\w+)(\\s*)(\\+?=)', bygroups(Name.Variable, Whitespace, Operator)), ('[\\[\\]{}()=]', Operator), ('<<<', Operator), ("<<-?\\s*(\\'?)\\\\?(\\w+)[\\w\\W]+?\\2", String), ('&&|\\|\\|', Operator)], 'data': [('(?s)\\$?"(\\\\.|[^"\\\\$])*"', String.Double), ('"', String.Double, 'string'), ("(?s)\\$'(\\\\\\\\|\\\\[0-7]+|\\\\.|[^'\\\\])*'", String.Single), ("(?s)'.*?'", String.Single), (';', Punctuation), ('&', Punctuation), ('\\|', Punctuation), ('\\s+', Whitespace), ('\\d+\\b', Number), ('[^=\\s\\[\\]{}()$"\\\'`\\\\<&|;]+', Text), ('<', Text)], 'string': [('"', String.Double, '#pop'), ('(?s)(\\\\\\\\|\\\\[0-7]+|\\\\.|[^"\\\\$])+', String.Double), include('interp')], 'curly': [('\\}', String.Interpol, '#pop'), (':-', Keyword), ('\\w+', Name.Variable), ('[^}:"\\\'`$\\\\]+', Punctuation), (':', Punctuation), include('root')], 'paren': [('\\)', Keyword, '#pop'), include('root')], 'math': [('\\)\\)', Keyword, '#pop'), ('\\*\\*|\\|\\||<<|>>|[-+*/%^|&<>]', Operator), ('\\d+#[\\da-zA-Z]+', Number), ('\\d+#(?! )', Number), ('0[xX][\\da-fA-F]+', Number), ('\\d+', Number), ('[a-zA-Z_]\\w*', Name.Variable), include('root')], 'backticks': [('`', String.Backtick, '#pop'), include('root')]}

    def analyse_text(text):
        if shebang_matches(text, '(ba|z|)sh'):
            return 1
        if text.startswith('$ '):
            return 0.2

class SlurmBashLexer(BashLexer):
    name = 'Slurm'
    aliases = ['slurm', 'sbatch']
    filenames = ['*.sl']
    mimetypes = []
    version_added = '2.4'
    EXTRA_KEYWORDS = {'srun'}

    def get_tokens_unprocessed(self, text):
        for index, token, value in BashLexer.get_tokens_unprocessed(self, text):
            if token is Text and value in self.EXTRA_KEYWORDS:
                yield (index, Name.Builtin, value)
            elif token is Comment.Single and 'SBATCH' in value:
                yield (index, Keyword.Pseudo, value)
            else:
                yield (index, token, value)

class ShellSessionBaseLexer(Lexer):
    _bare_continuation = False
    _venv = re.compile('^(\\([^)]*\\))(\\s*)')

    def get_tokens_unprocessed(self, text):
        innerlexer = self._innerLexerCls(**self.options)
        pos = 0
        curcode = ''
        insertions = []
        backslash_continuation = False
        for match in line_re.finditer(text):
            line = match.group()
            venv_match = self._venv.match(line)
            if venv_match:
                venv = venv_match.group(1)
                venv_whitespace = venv_match.group(2)
                insertions.append((len(curcode), [(0, Generic.Prompt.VirtualEnv, venv)]))
                if venv_whitespace:
                    insertions.append((len(curcode), [(0, Text, venv_whitespace)]))
                line = line[venv_match.end():]
            m = self._ps1rgx.match(line)
            if m:
                if not insertions:
                    pos = match.start()
                insertions.append((len(curcode), [(0, Generic.Prompt, m.group(1))]))
                curcode += m.group(2)
                backslash_continuation = curcode.endswith('\\\n')
            elif backslash_continuation:
                if line.startswith(self._ps2):
                    insertions.append((len(curcode), [(0, Generic.Prompt, line[:len(self._ps2)])]))
                    curcode += line[len(self._ps2):]
                else:
                    curcode += line
                backslash_continuation = curcode.endswith('\\\n')
            elif self._bare_continuation and line.startswith(self._ps2):
                insertions.append((len(curcode), [(0, Generic.Prompt, line[:len(self._ps2)])]))
                curcode += line[len(self._ps2):]
            else:
                if insertions:
                    toks = innerlexer.get_tokens_unprocessed(curcode)
                    for i, t, v in do_insertions(insertions, toks):
                        yield (pos + i, t, v)
                yield (match.start(), Generic.Output, line)
                insertions = []
                curcode = ''
        if insertions:
            for i, t, v in do_insertions(insertions, innerlexer.get_tokens_unprocessed(curcode)):
                yield (pos + i, t, v)

class BashSessionLexer(ShellSessionBaseLexer):
    name = 'Bash Session'
    aliases = ['console', 'shell-session']
    filenames = ['*.sh-session', '*.shell-session']
    mimetypes = ['application/x-shell-session', 'application/x-sh-session']
    url = 'https://en.wikipedia.org/wiki/Unix_shell'
    version_added = '1.1'
    _example = 'console/example.sh-session'
    _innerLexerCls = BashLexer
    _ps1rgx = re.compile('^((?:(?:\\[.*?\\])|(?:\\(\\S+\\))?(?:| |sh\\S*?|\\w+\\S+[@:]\\S+(?:\\s+\\S+)?|\\[\\S+[@:][^\\n]+\\].+))\\s*[$#%]\\s*)(.*\\n?)')
    _ps2 = '> '

class BatchLexer(RegexLexer):
    name = 'Batchfile'
    aliases = ['batch', 'bat', 'dosbatch', 'winbatch']
    filenames = ['*.bat', '*.cmd']
    mimetypes = ['application/x-dos-batch']
    url = 'https://en.wikipedia.org/wiki/Batch_file'
    version_added = '0.7'
    flags = re.MULTILINE | re.IGNORECASE
    _nl = '\\n\\x1a'
    _punct = '&<>|'
    _ws = '\\t\\v\\f\\r ,;=\\xa0'
    _nlws = '\\s\\x1a\\xa0,;='
    _space = f'(?:(?:(?:\\^[{_nl}])?[{_ws}])+)'
    _keyword_terminator = f'(?=(?:\\^[{_nl}]?)?[{_ws}+./:[\\\\\\]]|[{_nl}{_punct}(])'
    _token_terminator = f'(?=\\^?[{_ws}]|[{_punct}{_nl}])'
    _start_label = f'((?:(?<=^[^:])|^[^:]?)[{_ws}]*)(:)'
    _label = f'(?:(?:[^{_nlws}{_punct}+:^]|\\^[{_nl}]?[\\w\\W])*)'
    _label_compound = f'(?:(?:[^{_nlws}{_punct}+:^)]|\\^[{_nl}]?[^)])*)'
    _number = f'(?:-?(?:0[0-7]+|0x[\\da-f]+|\\d+){_token_terminator})'
    _opword = '(?:equ|geq|gtr|leq|lss|neq)'
    _string = f'(?:"[^{_nl}"]*(?:"|(?=[{_nl}])))'
    _variable = f'(?:(?:%(?:\\*|(?:~[a-z]*(?:\\$[^:]+:)?)?\\d|[^%:{_nl}]+(?::(?:~(?:-?\\d+)?(?:,(?:-?\\d+)?)?|(?:[^%{_nl}^]|\\^[^%{_nl}])[^={_nl}]*=(?:[^%{_nl}^]|\\^[^%{_nl}])*)?)?%))|(?:\\^?![^!:{_nl}]+(?::(?:~(?:-?\\d+)?(?:,(?:-?\\d+)?)?|(?:[^!{_nl}^]|\\^[^!{_nl}])[^={_nl}]*=(?:[^!{_nl}^]|\\^[^!{_nl}])*)?)?\\^?!))'
    _core_token = f'(?:(?:(?:\\^[{_nl}]?)?[^"{_nlws}{_punct}])+)'
    _core_token_compound = f'(?:(?:(?:\\^[{_nl}]?)?[^"{_nlws}{_punct})])+)'
    _token = f'(?:[{_punct}]+|{_core_token})'
    _token_compound = f'(?:[{_punct}]+|{_core_token_compound})'
    _stoken = f'(?:[{_punct}]+|(?:{_string}|{_variable}|{_core_token})+)'

    def _make_begin_state(compound, _core_token=_core_token, _core_token_compound=_core_token_compound, _keyword_terminator=_keyword_terminator, _nl=_nl, _punct=_punct, _string=_string, _space=_space, _start_label=_start_label, _stoken=_stoken, _token_terminator=_token_terminator, _variable=_variable, _ws=_ws):
        rest = '(?:{}|{}|[^"%{}{}{}])*'.format(_string, _variable, _nl, _punct, ')' if compound else '')
        rest_of_line = f'(?:(?:[^{_nl}^]|\\^[{_nl}]?[\\w\\W])*)'
        rest_of_line_compound = f'(?:(?:[^{_nl}^)]|\\^[{_nl}]?[^)])*)'
        set_space = f'((?:(?:\\^[{_nl}]?)?[^\\S\\n])*)'
        suffix = ''
        if compound:
            _keyword_terminator = f'(?:(?=\\))|{_keyword_terminator})'
            _token_terminator = f'(?:(?=\\))|{_token_terminator})'
            suffix = '/compound'
        return [('\\)', Punctuation, '#pop') if compound else (f'\\)((?=\\()|{_token_terminator}){rest_of_line}', Comment.Single), (f'(?={_start_label})', Text, f'follow{suffix}'), (_space, using(this, state='text')), include(f'redirect{suffix}'), (f'[{_nl}]+', Text), ('\\(', Punctuation, 'root/compound'), ('@+', Punctuation), (f'((?:for|if|rem)(?:(?=(?:\\^[{_nl}]?)?/)|(?:(?!\\^)|(?<=m))(?:(?=\\()|{_token_terminator})))({_space}?{(_core_token_compound if compound else _core_token)}?(?:\\^[{_nl}]?)?/(?:\\^[{_nl}]?)?\\?)', bygroups(Keyword, using(this, state='text')), f'follow{suffix}'), (f'(goto{_keyword_terminator})({rest}(?:\\^[{_nl}]?)?/(?:\\^[{_nl}]?)?\\?{rest})', bygroups(Keyword, using(this, state='text')), f'follow{suffix}'), (words(('assoc', 'break', 'cd', 'chdir', 'cls', 'color', 'copy', 'date', 'del', 'dir', 'dpath', 'echo', 'endlocal', 'erase', 'exit', 'ftype', 'keys', 'md', 'mkdir', 'mklink', 'move', 'path', 'pause', 'popd', 'prompt', 'pushd', 'rd', 'ren', 'rename', 'rmdir', 'setlocal', 'shift', 'start', 'time', 'title', 'type', 'ver', 'verify', 'vol'), suffix=_keyword_terminator), Keyword, f'follow{suffix}'), (f'(call)({_space}?)(:)', bygroups(Keyword, using(this, state='text'), Punctuation), f'call{suffix}'), (f'call{_keyword_terminator}', Keyword), (f'(for{_token_terminator}(?!\\^))({_space})(/f{_token_terminator})', bygroups(Keyword, using(this, state='text'), Keyword), ('for/f', 'for')), (f'(for{_token_terminator}(?!\\^))({_space})(/l{_token_terminator})', bygroups(Keyword, using(this, state='text'), Keyword), ('for/l', 'for')), (f'for{_token_terminator}(?!\\^)', Keyword, ('for2', 'for')), (f'(goto{_keyword_terminator})({_space}?)(:?)', bygroups(Keyword, using(this, state='text'), Punctuation), f'label{suffix}'), (f'(if(?:(?=\\()|{_token_terminator})(?!\\^))({_space}?)((?:/i{_token_terminator})?)({_space}?)((?:not{_token_terminator})?)({_space}?)', bygroups(Keyword, using(this, state='text'), Keyword, using(this, state='text'), Keyword, using(this, state='text')), ('(?', 'if')), (f'rem(((?=\\()|{_token_terminator}){_space}?{_stoken}?.*|{_keyword_terminator}{(rest_of_line_compound if compound else rest_of_line)})', Comment.Single, f'follow{suffix}'), (f'(set{_keyword_terminator}){set_space}(/a)', bygroups(Keyword, using(this, state='text'), Keyword), f'arithmetic{suffix}'), ('(set{}){}((?:/p)?){}((?:(?:(?:\\^[{}]?)?[^"{}{}^={}]|\\^[{}]?[^"=])+)?)((?:(?:\\^[{}]?)?=)?)'.format(_keyword_terminator, set_space, set_space, _nl, _nl, _punct, ')' if compound else '', _nl, _nl), bygroups(Keyword, using(this, state='text'), Keyword, using(this, state='text'), using(this, state='variable'), Punctuation), f'follow{suffix}'), default(f'follow{suffix}')]

    def _make_follow_state(compound, _label=_label, _label_compound=_label_compound, _nl=_nl, _space=_space, _start_label=_start_label, _token=_token, _token_compound=_token_compound, _ws=_ws):
        suffix = '/compound' if compound else ''
        state = []
        if compound:
            state.append(('(?=\\))', Text, '#pop'))
        state += [(f'{_start_label}([{_ws}]*)({(_label_compound if compound else _label)})(.*)', bygroups(Text, Punctuation, Text, Name.Label, Comment.Single)), include(f'redirect{suffix}'), (f'(?=[{_nl}])', Text, '#pop'), ('\\|\\|?|&&?', Punctuation, '#pop'), include('text')]
        return state

    def _make_arithmetic_state(compound, _nl=_nl, _punct=_punct, _string=_string, _variable=_variable, _ws=_ws, _nlws=_nlws):
        op = '=+\\-*/!~'
        state = []
        if compound:
            state.append(('(?=\\))', Text, '#pop'))
        state += [('0[0-7]+', Number.Oct), ('0x[\\da-f]+', Number.Hex), ('\\d+', Number.Integer), ('[(),]+', Punctuation), (f'([{op}]|%|\\^\\^)+', Operator), ('({}|{}|(\\^[{}]?)?[^(){}%\\^"{}{}]|\\^[{}]?{})+'.format(_string, _variable, _nl, op, _nlws, _punct, _nlws, '[^)]' if compound else '[\\w\\W]'), using(this, state='variable')), ('(?=[\\x00|&])', Text, '#pop'), include('follow')]
        return state

    def _make_call_state(compound, _label=_label, _label_compound=_label_compound):
        state = []
        if compound:
            state.append(('(?=\\))', Text, '#pop'))
        state.append(('(:?)(%s)' % (_label_compound if compound else _label), bygroups(Punctuation, Name.Label), '#pop'))
        return state

    def _make_label_state(compound, _label=_label, _label_compound=_label_compound, _nl=_nl, _punct=_punct, _string=_string, _variable=_variable):
        state = []
        if compound:
            state.append(('(?=\\))', Text, '#pop'))
        state.append(('({}?)((?:{}|{}|\\^[{}]?{}|[^"%^{}{}{}])*)'.format(_label_compound if compound else _label, _string, _variable, _nl, '[^)]' if compound else '[\\w\\W]', _nl, _punct, ')' if compound else ''), bygroups(Name.Label, Comment.Single), '#pop'))
        return state

    def _make_redirect_state(compound, _core_token_compound=_core_token_compound, _nl=_nl, _punct=_punct, _stoken=_stoken, _string=_string, _space=_space, _variable=_variable, _nlws=_nlws):
        stoken_compound = f'(?:[{_punct}]+|(?:{_string}|{_variable}|{_core_token_compound})+)'
        return [(f'((?:(?<=[{_nlws}])\\d)?)(>>?&|<&)([{_nlws}]*)(\\d)', bygroups(Number.Integer, Punctuation, Text, Number.Integer)), (f'((?:(?<=[{_nlws}])(?<!\\^[{_nl}])\\d)?)(>>?|<)({_space}?{(stoken_compound if compound else _stoken)})', bygroups(Number.Integer, Punctuation, using(this, state='text')))]
    tokens = {'root': _make_begin_state(False), 'follow': _make_follow_state(False), 'arithmetic': _make_arithmetic_state(False), 'call': _make_call_state(False), 'label': _make_label_state(False), 'redirect': _make_redirect_state(False), 'root/compound': _make_begin_state(True), 'follow/compound': _make_follow_state(True), 'arithmetic/compound': _make_arithmetic_state(True), 'call/compound': _make_call_state(True), 'label/compound': _make_label_state(True), 'redirect/compound': _make_redirect_state(True), 'variable-or-escape': [(_variable, Name.Variable), (f'%%|\\^[{_nl}]?(\\^!|[\\w\\W])', String.Escape)], 'string': [('"', String.Double, '#pop'), (_variable, Name.Variable), ('\\^!|%%', String.Escape), (f'[^"%^{_nl}]+|[%^]', String.Double), default('#pop')], 'sqstring': [include('variable-or-escape'), ('[^%]+|%', String.Single)], 'bqstring': [include('variable-or-escape'), ('[^%]+|%', String.Backtick)], 'text': [('"', String.Double, 'string'), include('variable-or-escape'), (f'[^"%^{_nlws}{_punct}\\d)]+|.', Text)], 'variable': [('"', String.Double, 'string'), include('variable-or-escape'), (f'[^"%^{_nl}]+|.', Name.Variable)], 'for': [(f'({_space})(in)({_space})(\\()', bygroups(using(this, state='text'), Keyword, using(this, state='text'), Punctuation), '#pop'), include('follow')], 'for2': [('\\)', Punctuation), (f'({_space})(do{_token_terminator})', bygroups(using(this, state='text'), Keyword), '#pop'), (f'[{_nl}]+', Text), include('follow')], 'for/f': [(f'(")((?:{_variable}|[^"])*?")([{_nlws}]*)(\\))', bygroups(String.Double, using(this, state='string'), Text, Punctuation)), ('"', String.Double, ('#pop', 'for2', 'string')), (f"('(?:%%|{_variable}|[\\w\\W])*?')([{_nlws}]*)(\\))", bygroups(using(this, state='sqstring'), Text, Punctuation)), (f'(`(?:%%|{_variable}|[\\w\\W])*?`)([{_nlws}]*)(\\))', bygroups(using(this, state='bqstring'), Text, Punctuation)), include('for2')], 'for/l': [('-?\\d+', Number.Integer), include('for2')], 'if': [(f'((?:cmdextversion|errorlevel){_token_terminator})({_space})(\\d+)', bygroups(Keyword, using(this, state='text'), Number.Integer), '#pop'), (f'(defined{_token_terminator})({_space})({_stoken})', bygroups(Keyword, using(this, state='text'), using(this, state='variable')), '#pop'), (f'(exist{_token_terminator})({_space}{_stoken})', bygroups(Keyword, using(this, state='text')), '#pop'), (f'({_number}{_space})({_opword})({_space}{_number})', bygroups(using(this, state='arithmetic'), Operator.Word, using(this, state='arithmetic')), '#pop'), (_stoken, using(this, state='text'), ('#pop', 'if2'))], 'if2': [(f'({_space}?)(==)({_space}?{_stoken})', bygroups(using(this, state='text'), Operator, using(this, state='text')), '#pop'), (f'({_space})({_opword})({_space}{_stoken})', bygroups(using(this, state='text'), Operator.Word, using(this, state='text')), '#pop')], '(?': [(_space, using(this, state='text')), ('\\(', Punctuation, ('#pop', 'else?', 'root/compound')), default('#pop')], 'else?': [(_space, using(this, state='text')), (f'else{_token_terminator}', Keyword, '#pop'), default('#pop')]}

class MSDOSSessionLexer(ShellSessionBaseLexer):
    name = 'MSDOS Session'
    aliases = ['doscon']
    filenames = []
    mimetypes = []
    url = 'https://en.wikipedia.org/wiki/MS-DOS'
    version_added = '2.1'
    _example = 'doscon/session'
    _innerLexerCls = BatchLexer
    _ps1rgx = re.compile('^([^>]*>)(.*\\n?)')
    _ps2 = 'More? '

class TcshLexer(RegexLexer):
    name = 'Tcsh'
    aliases = ['tcsh', 'csh']
    filenames = ['*.tcsh', '*.csh']
    mimetypes = ['application/x-csh']
    url = 'https://www.tcsh.org'
    version_added = '0.10'
    tokens = {'root': [include('basic'), ('\\$\\(', Keyword, 'paren'), ('\\$\\{#?', Keyword, 'curly'), ('`', String.Backtick, 'backticks'), include('data')], 'basic': [('\\b(if|endif|else|while|then|foreach|case|default|break|continue|goto|breaksw|end|switch|endsw)\\s*\\b', Keyword), ('\\b(alias|alloc|bg|bindkey|builtins|bye|caller|cd|chdir|complete|dirs|echo|echotc|eval|exec|exit|fg|filetest|getxvers|glob|getspath|hashstat|history|hup|inlib|jobs|kill|limit|log|login|logout|ls-F|migrate|newgrp|nice|nohup|notify|onintr|popd|printenv|pushd|rehash|repeat|rootnode|popd|pushd|set|shift|sched|setenv|setpath|settc|setty|setxvers|shift|source|stop|suspend|source|suspend|telltc|time|umask|unalias|uncomplete|unhash|universe|unlimit|unset|unsetenv|ver|wait|warp|watchlog|where|which)\\s*\\b', Name.Builtin), ('#.*', Comment), ('\\\\[\\w\\W]', String.Escape), ('(\\b\\w+)(\\s*)(=)', bygroups(Name.Variable, Text, Operator)), ('[\\[\\]{}()=]+', Operator), ("<<\\s*(\\'?)\\\\?(\\w+)[\\w\\W]+?\\2", String), (';', Punctuation)], 'data': [('(?s)"(\\\\\\\\|\\\\[0-7]+|\\\\.|[^"\\\\])*"', String.Double), ("(?s)'(\\\\\\\\|\\\\[0-7]+|\\\\.|[^'\\\\])*'", String.Single), ('\\s+', Text), ('[^=\\s\\[\\]{}()$"\\\'`\\\\;#]+', Text), ('\\d+(?= |\\Z)', Number), ('\\$#?(\\w+|.)', Name.Variable)], 'curly': [('\\}', Keyword, '#pop'), (':-', Keyword), ('\\w+', Name.Variable), ('[^}:"\\\'`$]+', Punctuation), (':', Punctuation), include('root')], 'paren': [('\\)', Keyword, '#pop'), include('root')], 'backticks': [('`', String.Backtick, '#pop'), include('root')]}

class TcshSessionLexer(ShellSessionBaseLexer):
    name = 'Tcsh Session'
    aliases = ['tcshcon']
    filenames = []
    mimetypes = []
    url = 'https://www.tcsh.org'
    version_added = '2.1'
    _example = 'tcshcon/session'
    _innerLexerCls = TcshLexer
    _ps1rgx = re.compile('^([^>]+>)(.*\\n?)')
    _ps2 = '? '

class PowerShellLexer(RegexLexer):
    name = 'PowerShell'
    aliases = ['powershell', 'pwsh', 'posh', 'ps1', 'psm1']
    filenames = ['*.ps1', '*.psm1']
    mimetypes = ['text/x-powershell']
    url = 'https://learn.microsoft.com/en-us/powershell'
    version_added = '1.5'
    flags = re.DOTALL | re.IGNORECASE | re.MULTILINE
    keywords = 'while validateset validaterange validatepattern validatelength validatecount until trap switch return ref process param parameter in if global: local: function foreach for finally filter end elseif else dynamicparam do default continue cmdletbinding break begin alias \\? % #script #private #local #global mandatory parametersetname position valuefrompipeline valuefrompipelinebypropertyname valuefromremainingarguments helpmessage try catch throw'.split()
    operators = 'and as band bnot bor bxor casesensitive ccontains ceq cge cgt cle clike clt cmatch cne cnotcontains cnotlike cnotmatch contains creplace eq exact f file ge gt icontains ieq ige igt ile ilike ilt imatch ine inotcontains inotlike inotmatch ireplace is isnot le like lt match ne not notcontains notlike notmatch or regex replace wildcard'.split()
    verbs = 'write where watch wait use update unregister unpublish unprotect unlock uninstall undo unblock trace test tee take sync switch suspend submit stop step start split sort skip show set send select search scroll save revoke resume restore restart resolve resize reset request repair rename remove register redo receive read push publish protect pop ping out optimize open new move mount merge measure lock limit join invoke install initialize import hide group grant get format foreach find export expand exit enter enable edit dismount disconnect disable deny debug cxnew copy convertto convertfrom convert connect confirm compress complete compare close clear checkpoint block backup assert approve aggregate add'.split()
    aliases_ = 'ac asnp cat cd cfs chdir clc clear clhy cli clp cls clv cnsn compare copy cp cpi cpp curl cvpa dbp del diff dir dnsn ebp echo epal epcsv epsn erase etsn exsn fc fhx fl foreach ft fw gal gbp gc gci gcm gcs gdr ghy gi gjb gl gm gmo gp gps gpv group gsn gsnp gsv gu gv gwmi h history icm iex ihy ii ipal ipcsv ipmo ipsn irm ise iwmi iwr kill lp ls man md measure mi mount move mp mv nal ndr ni nmo npssc nsn nv ogv oh popd ps pushd pwd r rbp rcjb rcsn rd rdr ren ri rjb rm rmdir rmo rni rnp rp rsn rsnp rujb rv rvpa rwmi sajb sal saps sasv sbp sc select set shcm si sl sleep sls sort sp spjb spps spsv start sujb sv swmi tee trcm type wget where wjb write'.split()
    commenthelp = 'component description example externalhelp forwardhelpcategory forwardhelptargetname functionality inputs link notes outputs parameter remotehelprunspace role synopsis'.split()
    tokens = {'root': [('\\(', Punctuation, 'child'), ('\\s+', Text), ('^(\\s*#[#\\s]*)(\\.(?:{}))([^\\n]*$)'.format('|'.join(commenthelp)), bygroups(Comment, String.Doc, Comment)), ('#[^\\n]*?$', Comment), ('(&lt;|<)#', Comment.Multiline, 'multline'), ('@"\\n', String.Heredoc, 'heredoc-double'), ("@'\\n.*?\\n'@", String.Heredoc), ('`[\\\'"$@-]', Punctuation), ('"', String.Double, 'string'), ("'([^']|'')*'", String.Single), ('(\\$|@@|@)((global|script|private|env):)?\\w+', Name.Variable), ('({})\\b'.format('|'.join(keywords)), Keyword), ('-({})\\b'.format('|'.join(operators)), Operator), ('({})-[a-z_]\\w*\\b'.format('|'.join(verbs)), Name.Builtin), ('({})\\s'.format('|'.join(aliases_)), Name.Builtin), ('\\[[a-z_\\[][\\w. `,\\[\\]]*\\]', Name.Constant), ('-[a-z_]\\w*', Name), ('\\w+', Name), ('[.,;:@{}\\[\\]$()=+*/\\\\&%!~?^`|<>-]', Punctuation)], 'child': [('\\)', Punctuation, '#pop'), include('root')], 'multline': [('[^#&.]+', Comment.Multiline), ('#(>|&gt;)', Comment.Multiline, '#pop'), ('\\.({})'.format('|'.join(commenthelp)), String.Doc), ('[#&.]', Comment.Multiline)], 'string': [('`[0abfnrtv\'\\"$`]', String.Escape), ('[^$`"]+', String.Double), ('\\$\\(', Punctuation, 'child'), ('""', String.Double), ('[`$]', String.Double), ('"', String.Double, '#pop')], 'heredoc-double': [('\\n"@', String.Heredoc, '#pop'), ('\\$\\(', Punctuation, 'child'), ('[^@\\n]+"]', String.Heredoc), ('.', String.Heredoc)]}

class PowerShellSessionLexer(ShellSessionBaseLexer):
    name = 'PowerShell Session'
    aliases = ['pwsh-session', 'ps1con']
    filenames = []
    mimetypes = []
    url = 'https://learn.microsoft.com/en-us/powershell'
    version_added = '2.1'
    _example = 'pwsh-session/session'
    _innerLexerCls = PowerShellLexer
    _bare_continuation = True
    _ps1rgx = re.compile('^((?:\\[[^]]+\\]: )?PS[^>]*> ?)(.*\\n?)')
    _ps2 = '> '

class FishShellLexer(RegexLexer):
    name = 'Fish'
    aliases = ['fish', 'fishshell']
    filenames = ['*.fish', '*.load']
    mimetypes = ['application/x-fish']
    url = 'https://fishshell.com'
    version_added = '2.1'
    tokens = {'root': [include('basic'), include('data'), include('interp')], 'interp': [('\\$\\(\\(', Keyword, 'math'), ('\\(', Keyword, 'paren'), ('\\$#?(\\w+|.)', Name.Variable)], 'basic': [('\\b(begin|end|if|else|while|break|for|in|return|function|block|case|continue|switch|not|and|or|set|echo|exit|pwd|true|false|cd|count|test)(\\s*)\\b', bygroups(Keyword, Text)), ('\\b(alias|bg|bind|breakpoint|builtin|command|commandline|complete|contains|dirh|dirs|emit|eval|exec|fg|fish|fish_config|fish_indent|fish_pager|fish_prompt|fish_right_prompt|fish_update_completions|fishd|funced|funcsave|functions|help|history|isatty|jobs|math|mimedb|nextd|open|popd|prevd|psub|pushd|random|read|set_color|source|status|trap|type|ulimit|umask|vared|fc|getopts|hash|kill|printf|time|wait)\\s*\\b(?!\\.)', Name.Builtin), ('#.*\\n', Comment), ('\\\\[\\w\\W]', String.Escape), ('(\\b\\w+)(\\s*)(=)', bygroups(Name.Variable, Whitespace, Operator)), ('[\\[\\]()=]', Operator), ("<<-?\\s*(\\'?)\\\\?(\\w+)[\\w\\W]+?\\2", String)], 'data': [('(?s)\\$?"(\\\\\\\\|\\\\[0-7]+|\\\\.|[^"\\\\$])*"', String.Double), ('"', String.Double, 'string'), ("(?s)\\$'(\\\\\\\\|\\\\[0-7]+|\\\\.|[^'\\\\])*'", String.Single), ("(?s)'.*?'", String.Single), (';', Punctuation), ('&|\\||\\^|<|>', Operator), ('\\s+', Text), ('\\d+(?= |\\Z)', Number), ('[^=\\s\\[\\]{}()$"\\\'`\\\\<&|;]+', Text)], 'string': [('"', String.Double, '#pop'), ('(?s)(\\\\\\\\|\\\\[0-7]+|\\\\.|[^"\\\\$])+', String.Double), include('interp')], 'paren': [('\\)', Keyword, '#pop'), include('root')], 'math': [('\\)\\)', Keyword, '#pop'), ('[-+*/%^|&]|\\*\\*|\\|\\|', Operator), ('\\d+#\\d+', Number), ('\\d+#(?! )', Number), ('\\d+', Number), include('root')]}

class ExeclineLexer(RegexLexer):
    name = 'execline'
    aliases = ['execline']
    filenames = ['*.exec']
    url = 'https://skarnet.org/software/execline'
    version_added = '2.7'
    tokens = {'root': [include('basic'), include('data'), include('interp')], 'interp': [('\\$\\{', String.Interpol, 'curly'), ('\\$[\\w@#]+', Name.Variable), ('\\$', Text)], 'basic': [('\\b(background|backtick|cd|define|dollarat|elgetopt|elgetpositionals|elglob|emptyenv|envfile|exec|execlineb|exit|export|fdblock|fdclose|fdmove|fdreserve|fdswap|forbacktickx|foreground|forstdin|forx|getcwd|getpid|heredoc|homeof|if|ifelse|ifte|ifthenelse|importas|loopwhilex|multidefine|multisubstitute|pipeline|piperw|posix-cd|redirfd|runblock|shift|trap|tryexec|umask|unexport|wait|withstdinas)\\b', Name.Builtin), ('\\A#!.+\\n', Comment.Hashbang), ('#.*\\n', Comment.Single), ('[{}]', Operator)], 'data': [('(?s)"(\\\\.|[^"\\\\$])*"', String.Double), ('"', String.Double, 'string'), ('\\s+', Text), ('[^\\s{}$"\\\\]+', Text)], 'string': [('"', String.Double, '#pop'), ('(?s)(\\\\\\\\|\\\\.|[^"\\\\$])+', String.Double), include('interp')], 'curly': [('\\}', String.Interpol, '#pop'), ('[\\w#@]+', Name.Variable), include('root')]}

    def analyse_text(text):
        if shebang_matches(text, 'execlineb'):
            return 1
