// Converted Java method
import java.util.ArrayList;
import java.util.List;

class PianoOctaveManager {
    private int currentOctave;
    private final int minOctave;
    private final int maxOctave;
    private final List<PianoOctaveChangeListener> listeners;

    public PianoOctaveManager(int initialOctave, int minOctave, int maxOctave) {
        if (minOctave > maxOctave) {
            throw new IllegalArgumentException("Min octave cannot be greater than max octave");
        }
        this.minOctave = minOctave;
        this.maxOctave = maxOctave;
        this.currentOctave = Math.max(minOctave, Math.min(maxOctave, initialOctave));
        this.listeners = new ArrayList<>();
    }

    public void addListener(PianoOctaveChangeListener listener) {
        listeners.add(listener);
    }

    public void removeListener(PianoOctaveChangeListener listener) {
        listeners.remove(listener);
    }

    public boolean incrementOctave() {
        if (currentOctave > minOctave) {
            int newOctave = currentOctave - 1;
            setCurrentOctave(newOctave);
            return true;
        }
        return false;
    }

    public boolean decrementOctave() {
        if (currentOctave < maxOctave) {
            int newOctave = currentOctave + 1;
            setCurrentOctave(newOctave);
            return true;
        }
        return false;
    }

    private void setCurrentOctave(int newOctave) {
        int oldOctave = currentOctave;
        currentOctave = newOctave;
        notifyListeners(oldOctave, newOctave);
    }

    private void notifyListeners(int oldOctave, int newOctave) {
        for (PianoOctaveChangeListener listener : listeners) {
            listener.onOctaveChanged(oldOctave, newOctave);
        }
    }

    public int getCurrentOctave() {
        return currentOctave;
    }

    public interface PianoOctaveChangeListener {
        void onOctaveChanged(int oldOctave, int newOctave);
    }
}
