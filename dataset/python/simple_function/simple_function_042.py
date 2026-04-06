import re
import sys
from html import escape
from weakref import proxy
from .std import tqdm as std_tqdm
if True:
    IPY = 0
    try:
        import ipywidgets
        IPY = 4
    except ImportError:
        IPY = 32
        import warnings
        with warnings.catch_warnings():
            warnings.filterwarnings('ignore', message='.*The `IPython.html` package has been deprecated.*')
            try:
                import IPython.html.widgets as ipywidgets
            except ImportError:
                pass
    try:
        if IPY == 32:
            from IPython.html.widgets import HTML
            from IPython.html.widgets import FloatProgress as IProgress
            from IPython.html.widgets import HBox
            IPY = 3
        else:
            from ipywidgets import HTML
            from ipywidgets import FloatProgress as IProgress
            from ipywidgets import HBox
    except ImportError:
        try:
            from IPython.html.widgets import HTML
            from IPython.html.widgets import ContainerWidget as HBox
            from IPython.html.widgets import FloatProgressWidget as IProgress
            IPY = 2
        except ImportError:
            IPY = 0
            IProgress = None
            HBox = object
    try:
        from IPython.display import display
    except ImportError:
        pass
__author__ = {'github.com/': ['lrq3000', 'casperdcl', 'alexanderkuk']}
__all__ = ['tqdm_notebook', 'tnrange', 'tqdm', 'trange']
WARN_NOIPYW = 'IProgress not found. Please update jupyter and ipywidgets. See https://ipywidgets.readthedocs.io/en/stable/user_install.html'

class TqdmHBox(HBox):

    def _json_(self, pretty=None):
        pbar = getattr(self, 'pbar', None)
        if pbar is None:
            return {}
        d = pbar.format_dict
        if pretty is not None:
            d['ascii'] = not pretty
        return d

    def __repr__(self, pretty=False):
        pbar = getattr(self, 'pbar', None)
        if pbar is None:
            return super().__repr__()
        return pbar.format_meter(**self._json_(pretty))

    def _repr_pretty_(self, pp, *_, **__):
        pp.text(self.__repr__(True))

class tqdm_notebook(std_tqdm):

    @staticmethod
    def status_printer(_, total=None, desc=None, ncols=None):
        if IProgress is None:
            raise ImportError(WARN_NOIPYW)
        if total:
            pbar = IProgress(min=0, max=total)
        else:
            pbar = IProgress(min=0, max=1)
            pbar.value = 1
            pbar.bar_style = 'info'
            if ncols is None:
                pbar.layout.width = '20px'
        ltext = HTML()
        rtext = HTML()
        if desc:
            ltext.value = desc
        container = TqdmHBox(children=[ltext, pbar, rtext])
        if ncols is not None:
            ncols = str(ncols)
            try:
                if int(ncols) > 0:
                    ncols += 'px'
            except ValueError:
                pass
            pbar.layout.flex = '2'
            container.layout.width = ncols
            container.layout.display = 'inline-flex'
            container.layout.flex_flow = 'row wrap'
        return container

    def display(self, msg=None, pos=None, close=False, bar_style=None, check_delay=True):
        if not msg and (not close):
            d = self.format_dict
            d['bar_format'] = (d['bar_format'] or '{l_bar}<bar/>{r_bar}').replace('{bar}', '<bar/>')
            msg = self.format_meter(**d)
        ltext, pbar, rtext = self.container.children
        pbar.value = self.n
        if msg:
            msg = msg.replace(' ', '\u2007')
            if '<bar/>' in msg:
                left, right = map(escape, re.split('\\|?<bar/>\\|?', msg, maxsplit=1))
            else:
                left, right = ('', escape(msg))
            ltext.value = left
            if right:
                rtext.value = right
        if bar_style:
            if pbar.bar_style != 'danger' or bar_style != 'success':
                pbar.bar_style = bar_style
        if close and pbar.bar_style != 'danger':
            try:
                self.container.close()
            except AttributeError:
                self.container.visible = False
            self.container.layout.visibility = 'hidden'
        if check_delay and self.delay > 0 and (not self.displayed):
            display(self.container)
            self.displayed = True

    @property
    def colour(self):
        if hasattr(self, 'container'):
            return self.container.children[-2].style.bar_color

    @colour.setter
    def colour(self, bar_color):
        if hasattr(self, 'container'):
            self.container.children[-2].style.bar_color = bar_color

    def __init__(self, *args, **kwargs):
        kwargs = kwargs.copy()
        file_kwarg = kwargs.get('file', sys.stderr)
        if file_kwarg is sys.stderr or file_kwarg is None:
            kwargs['file'] = sys.stdout
        kwargs['gui'] = True
        kwargs['disable'] = bool(kwargs.get('disable', False))
        colour = kwargs.pop('colour', None)
        display_here = kwargs.pop('display', True)
        super().__init__(*args, **kwargs)
        if self.disable or not kwargs['gui']:
            self.disp = lambda *_, **__: None
            return
        self.ncols = '100%' if self.dynamic_ncols else kwargs.get('ncols', None)
        unit_scale = 1 if self.unit_scale is True else self.unit_scale or 1
        total = self.total * unit_scale if self.total else self.total
        self.container = self.status_printer(self.fp, total, self.desc, self.ncols)
        self.container.pbar = proxy(self)
        self.displayed = False
        if display_here and self.delay <= 0:
            display(self.container)
            self.displayed = True
        self.disp = self.display
        self.colour = colour
        if not self.disable:
            self.display(check_delay=False)

    def __iter__(self):
        try:
            it = super().__iter__()
            yield from it
        except:
            self.disp(bar_style='danger')
            raise

    def update(self, n=1):
        try:
            return super().update(n=n)
        except:
            self.disp(bar_style='danger')
            raise

    def close(self):
        if self.disable:
            return
        super().close()
        if self.total and self.n < self.total:
            self.disp(bar_style='danger', check_delay=False)
        elif self.leave:
            self.disp(bar_style='success', check_delay=False)
        else:
            self.disp(close=True, check_delay=False)

    def clear(self, *_, **__):
        pass

    def reset(self, total=None):
        if self.disable:
            return super().reset(total=total)
        _, pbar, _ = self.container.children
        pbar.bar_style = ''
        if total is not None:
            pbar.max = total
            if not self.total and self.ncols is None:
                pbar.layout.width = None
        return super().reset(total=total)

def tnrange(*args, **kwargs):
    return tqdm_notebook(range(*args), **kwargs)
tqdm = tqdm_notebook
trange = tnrange
