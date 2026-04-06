from io import StringIO
from pygments.formatter import Formatter
from pygments.lexer import Lexer, do_insertions
from pygments.token import Token, STANDARD_TYPES
from pygments.util import get_bool_opt, get_int_opt
__all__ = ['LatexFormatter']

def escape_tex(text, commandprefix):
    return text.replace('\\', '\x00').replace('{', '\x01').replace('}', '\x02').replace('\x00', f'\\{commandprefix}Zbs{{}}').replace('\x01', f'\\{commandprefix}Zob{{}}').replace('\x02', f'\\{commandprefix}Zcb{{}}').replace('^', f'\\{commandprefix}Zca{{}}').replace('_', f'\\{commandprefix}Zus{{}}').replace('&', f'\\{commandprefix}Zam{{}}').replace('<', f'\\{commandprefix}Zlt{{}}').replace('>', f'\\{commandprefix}Zgt{{}}').replace('#', f'\\{commandprefix}Zsh{{}}').replace('%', f'\\{commandprefix}Zpc{{}}').replace('$', f'\\{commandprefix}Zdl{{}}').replace('-', f'\\{commandprefix}Zhy{{}}').replace("'", f'\\{commandprefix}Zsq{{}}').replace('"', f'\\{commandprefix}Zdq{{}}').replace('~', f'\\{commandprefix}Zti{{}}')
DOC_TEMPLATE = '\n\\documentclass{%(docclass)s}\n\\usepackage{fancyvrb}\n\\usepackage{color}\n\\usepackage[%(encoding)s]{inputenc}\n%(preamble)s\n\n%(styledefs)s\n\n\\begin{document}\n\n\\section*{%(title)s}\n\n%(code)s\n\\end{document}\n'
STYLE_TEMPLATE = '\n\\makeatletter\n\\def\\%(cp)s@reset{\\let\\%(cp)s@it=\\relax \\let\\%(cp)s@bf=\\relax%%\n    \\let\\%(cp)s@ul=\\relax \\let\\%(cp)s@tc=\\relax%%\n    \\let\\%(cp)s@bc=\\relax \\let\\%(cp)s@ff=\\relax}\n\\def\\%(cp)s@tok#1{\\csname %(cp)s@tok@#1\\endcsname}\n\\def\\%(cp)s@toks#1+{\\ifx\\relax#1\\empty\\else%%\n    \\%(cp)s@tok{#1}\\expandafter\\%(cp)s@toks\\fi}\n\\def\\%(cp)s@do#1{\\%(cp)s@bc{\\%(cp)s@tc{\\%(cp)s@ul{%%\n    \\%(cp)s@it{\\%(cp)s@bf{\\%(cp)s@ff{#1}}}}}}}\n\\def\\%(cp)s#1#2{\\%(cp)s@reset\\%(cp)s@toks#1+\\relax+\\%(cp)s@do{#2}}\n\n%(styles)s\n\n\\def\\%(cp)sZbs{\\char`\\\\}\n\\def\\%(cp)sZus{\\char`\\_}\n\\def\\%(cp)sZob{\\char`\\{}\n\\def\\%(cp)sZcb{\\char`\\}}\n\\def\\%(cp)sZca{\\char`\\^}\n\\def\\%(cp)sZam{\\char`\\&}\n\\def\\%(cp)sZlt{\\char`\\<}\n\\def\\%(cp)sZgt{\\char`\\>}\n\\def\\%(cp)sZsh{\\char`\\#}\n\\def\\%(cp)sZpc{\\char`\\%%}\n\\def\\%(cp)sZdl{\\char`\\$}\n\\def\\%(cp)sZhy{\\char`\\-}\n\\def\\%(cp)sZsq{\\char`\\\'}\n\\def\\%(cp)sZdq{\\char`\\"}\n\\def\\%(cp)sZti{\\char`\\~}\n%% for compatibility with earlier versions\n\\def\\%(cp)sZat{@}\n\\def\\%(cp)sZlb{[}\n\\def\\%(cp)sZrb{]}\n\\makeatother\n'

def _get_ttype_name(ttype):
    fname = STANDARD_TYPES.get(ttype)
    if fname:
        return fname
    aname = ''
    while fname is None:
        aname = ttype[-1] + aname
        ttype = ttype.parent
        fname = STANDARD_TYPES.get(ttype)
    return fname + aname

class LatexFormatter(Formatter):
    name = 'LaTeX'
    aliases = ['latex', 'tex']
    filenames = ['*.tex']

    def __init__(self, **options):
        Formatter.__init__(self, **options)
        self.nowrap = get_bool_opt(options, 'nowrap', False)
        self.docclass = options.get('docclass', 'article')
        self.preamble = options.get('preamble', '')
        self.linenos = get_bool_opt(options, 'linenos', False)
        self.linenostart = abs(get_int_opt(options, 'linenostart', 1))
        self.linenostep = abs(get_int_opt(options, 'linenostep', 1))
        self.verboptions = options.get('verboptions', '')
        self.nobackground = get_bool_opt(options, 'nobackground', False)
        self.commandprefix = options.get('commandprefix', 'PY')
        self.texcomments = get_bool_opt(options, 'texcomments', False)
        self.mathescape = get_bool_opt(options, 'mathescape', False)
        self.escapeinside = options.get('escapeinside', '')
        if len(self.escapeinside) == 2:
            self.left = self.escapeinside[0]
            self.right = self.escapeinside[1]
        else:
            self.escapeinside = ''
        self.envname = options.get('envname', 'Verbatim')
        self._create_stylesheet()

    def _create_stylesheet(self):
        t2n = self.ttype2name = {Token: ''}
        c2d = self.cmd2def = {}
        cp = self.commandprefix

        def rgbcolor(col):
            if col:
                return ','.join(['%.2f' % (int(col[i] + col[i + 1], 16) / 255.0) for i in (0, 2, 4)])
            else:
                return '1,1,1'
        for ttype, ndef in self.style:
            name = _get_ttype_name(ttype)
            cmndef = ''
            if ndef['bold']:
                cmndef += '\\let\\$$@bf=\\textbf'
            if ndef['italic']:
                cmndef += '\\let\\$$@it=\\textit'
            if ndef['underline']:
                cmndef += '\\let\\$$@ul=\\underline'
            if ndef['roman']:
                cmndef += '\\let\\$$@ff=\\textrm'
            if ndef['sans']:
                cmndef += '\\let\\$$@ff=\\textsf'
            if ndef['mono']:
                cmndef += '\\let\\$$@ff=\\textsf'
            if ndef['color']:
                cmndef += '\\def\\$$@tc##1{{\\textcolor[rgb]{{{}}}{{##1}}}}'.format(rgbcolor(ndef['color']))
            if ndef['border']:
                cmndef += '\\def\\$$@bc##1{{{{\\setlength{{\\fboxsep}}{{\\string -\\fboxrule}}\\fcolorbox[rgb]{{{}}}{{{}}}{{\\strut ##1}}}}}}'.format(rgbcolor(ndef['border']), rgbcolor(ndef['bgcolor']))
            elif ndef['bgcolor']:
                cmndef += '\\def\\$$@bc##1{{{{\\setlength{{\\fboxsep}}{{0pt}}\\colorbox[rgb]{{{}}}{{\\strut ##1}}}}}}'.format(rgbcolor(ndef['bgcolor']))
            if cmndef == '':
                continue
            cmndef = cmndef.replace('$$', cp)
            t2n[ttype] = name
            c2d[name] = cmndef

    def get_style_defs(self, arg=''):
        cp = self.commandprefix
        styles = []
        for name, definition in self.cmd2def.items():
            styles.append(f'\\@namedef{{{cp}@tok@{name}}}{{{definition}}}')
        return STYLE_TEMPLATE % {'cp': self.commandprefix, 'styles': '\n'.join(styles)}

    def format_unencoded(self, tokensource, outfile):
        t2n = self.ttype2name
        cp = self.commandprefix
        if self.full:
            realoutfile = outfile
            outfile = StringIO()
        if not self.nowrap:
            outfile.write('\\begin{' + self.envname + '}[commandchars=\\\\\\{\\}')
            if self.linenos:
                start, step = (self.linenostart, self.linenostep)
                outfile.write(',numbers=left' + (start and ',firstnumber=%d' % start or '') + (step and ',stepnumber=%d' % step or ''))
            if self.mathescape or self.texcomments or self.escapeinside:
                outfile.write(',codes={\\catcode`\\$=3\\catcode`\\^=7\\catcode`\\_=8\\relax}')
            if self.verboptions:
                outfile.write(',' + self.verboptions)
            outfile.write(']\n')
        for ttype, value in tokensource:
            if ttype in Token.Comment:
                if self.texcomments:
                    start = value[0:1]
                    for i in range(1, len(value)):
                        if start[0] != value[i]:
                            break
                        start += value[i]
                    value = value[len(start):]
                    start = escape_tex(start, cp)
                    value = start + value
                elif self.mathescape:
                    parts = value.split('$')
                    in_math = False
                    for i, part in enumerate(parts):
                        if not in_math:
                            parts[i] = escape_tex(part, cp)
                        in_math = not in_math
                    value = '$'.join(parts)
                elif self.escapeinside:
                    text = value
                    value = ''
                    while text:
                        a, sep1, text = text.partition(self.left)
                        if sep1:
                            b, sep2, text = text.partition(self.right)
                            if sep2:
                                value += escape_tex(a, cp) + b
                            else:
                                value += escape_tex(a + sep1 + b, cp)
                        else:
                            value += escape_tex(a, cp)
                else:
                    value = escape_tex(value, cp)
            elif ttype not in Token.Escape:
                value = escape_tex(value, cp)
            styles = []
            while ttype is not Token:
                try:
                    styles.append(t2n[ttype])
                except KeyError:
                    styles.append(_get_ttype_name(ttype))
                ttype = ttype.parent
            styleval = '+'.join(reversed(styles))
            if styleval:
                spl = value.split('\n')
                for line in spl[:-1]:
                    if line:
                        outfile.write(f'\\{cp}{{{styleval}}}{{{line}}}')
                    outfile.write('\n')
                if spl[-1]:
                    outfile.write(f'\\{cp}{{{styleval}}}{{{spl[-1]}}}')
            else:
                outfile.write(value)
        if not self.nowrap:
            outfile.write('\\end{' + self.envname + '}\n')
        if self.full:
            encoding = self.encoding or 'utf8'
            encoding = {'utf_8': 'utf8', 'latin_1': 'latin1', 'iso_8859_1': 'latin1'}.get(encoding.replace('-', '_'), encoding)
            realoutfile.write(DOC_TEMPLATE % dict(docclass=self.docclass, preamble=self.preamble, title=self.title, encoding=encoding, styledefs=self.get_style_defs(), code=outfile.getvalue()))

class LatexEmbeddedLexer(Lexer):

    def __init__(self, left, right, lang, **options):
        self.left = left
        self.right = right
        self.lang = lang
        Lexer.__init__(self, **options)

    def get_tokens_unprocessed(self, text):
        buffered = ''
        insertions = []
        insertion_buf = []
        for i, t, v in self._find_safe_escape_tokens(text):
            if t is None:
                if insertion_buf:
                    insertions.append((len(buffered), insertion_buf))
                    insertion_buf = []
                buffered += v
            else:
                insertion_buf.append((i, t, v))
        if insertion_buf:
            insertions.append((len(buffered), insertion_buf))
        return do_insertions(insertions, self.lang.get_tokens_unprocessed(buffered))

    def _find_safe_escape_tokens(self, text):
        for i, t, v in self._filter_to(self.lang.get_tokens_unprocessed(text), lambda t: t in Token.Comment or t in Token.String):
            if t is None:
                for i2, t2, v2 in self._find_escape_tokens(v):
                    yield (i + i2, t2, v2)
            else:
                yield (i, None, v)

    def _filter_to(self, it, pred):
        buf = ''
        idx = 0
        for i, t, v in it:
            if pred(t):
                if buf:
                    yield (idx, None, buf)
                    buf = ''
                yield (i, t, v)
            else:
                if not buf:
                    idx = i
                buf += v
        if buf:
            yield (idx, None, buf)

    def _find_escape_tokens(self, text):
        index = 0
        while text:
            a, sep1, text = text.partition(self.left)
            if a:
                yield (index, None, a)
                index += len(a)
            if sep1:
                b, sep2, text = text.partition(self.right)
                if sep2:
                    yield (index + len(sep1), Token.Escape, b)
                    index += len(sep1) + len(b) + len(sep2)
                else:
                    yield (index, Token.Error, sep1)
                    index += len(sep1)
                    text = b
