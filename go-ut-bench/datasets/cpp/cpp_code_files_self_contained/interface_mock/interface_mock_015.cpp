#include <string>
#include <chrono>
#include <map>
#include <vector>
#include <functional>

using namespace std;
using namespace std::chrono;

class GameTimer {
private:
    system_clock::time_point start_time;
    milliseconds duration;
    bool is_running = false;
    bool is_paused = false;
    system_clock::time_point pause_time;
    milliseconds remaining_before_pause;

public:
    GameTimer() : duration(0) {}

    void start(milliseconds initial_duration) {
        duration = initial_duration;
        start_time = system_clock::now();
        is_running = true;
        is_paused = false;
    }

    void pause() {
        if (is_running && !is_paused) {
            pause_time = system_clock::now();
            remaining_before_pause = duration - duration_cast<milliseconds>(pause_time - start_time);
            is_paused = true;
        }
    }

    void resume() {
        if (is_running && is_paused) {
            start_time = system_clock::now() - (duration - remaining_before_pause);
            is_paused = false;
        }
    }

    void reset() {
        is_running = false;
        is_paused = false;
        duration = milliseconds(0);
    }

    void add_time(milliseconds additional_time) {
        if (is_paused) {
            remaining_before_pause += additional_time;
        } else if (is_running) {
            duration += additional_time;
        } else {
            duration = additional_time;
        }
    }

    milliseconds get_remaining_time() const {
        if (!is_running) return milliseconds(0);
        if (is_paused) return remaining_before_pause;
        auto elapsed = system_clock::now() - start_time;
        return duration - duration_cast<milliseconds>(elapsed);
    }

    string get_formatted_time() const {
        auto remaining = get_remaining_time();
        if (remaining <= milliseconds(0)) return "00:00";

        auto secs = duration_cast<seconds>(remaining);
        remaining -= duration_cast<milliseconds>(secs);
        auto mins = duration_cast<minutes>(secs);
        secs -= duration_cast<seconds>(mins);

        char buffer[6];
        snprintf(buffer, sizeof(buffer), "%02d:%02d", (int)mins.count(), (int)secs.count());
        return string(buffer);
    }

    bool is_active() const { return is_running && !is_paused; }
    bool is_expired() const { return get_remaining_time() <= milliseconds(0); }
};

class Team {
private:
    string name;
    int score = 0;
    int fouls = 0;
    int timeouts = 0;
    int max_timeouts;

public:
    Team(const string& team_name, int max_tos = 3) : name(team_name), max_timeouts(max_tos) {}

    void add_score(int points) { score += points; }
    void add_foul() { fouls++; }
    bool use_timeout() {
        if (timeouts < max_timeouts) {
            timeouts++;
            return true;
        }
        return false;
    }

    void reset_score() { score = 0; }
    void reset_fouls() { fouls = 0; }
    void reset_timeouts() { timeouts = 0; }

    string get_name() const { return name; }
    int get_score() const { return score; }
    int get_fouls() const { return fouls; }
    int get_timeouts() const { return timeouts; }
    int get_remaining_timeouts() const { return max_timeouts - timeouts; }
};

class AdvancedScoreboard {
private:
    Team home_team;
    Team away_team;
    GameTimer game_timer;
    GameTimer shot_clock;
    int quarter = 1;
    const int max_quarters = 4;
    bool is_game_over = false;
    vector<function<void()>> change_callbacks;

public:
    AdvancedScoreboard(const string& home_name, const string& away_name)
        : home_team(home_name), away_team(away_name), shot_clock() {
        shot_clock.start(24000ms); // 24-second shot clock
    }

    void register_callback(const function<void()>& callback) {
        change_callbacks.push_back(callback);
    }

    void notify_changes() {
        for (auto& callback : change_callbacks) {
            callback();
        }
    }

    void start_game(milliseconds game_duration) {
        game_timer.start(game_duration);
        is_game_over = false;
        notify_changes();
    }

    void pause_game() {
        game_timer.pause();
        shot_clock.pause();
        notify_changes();
    }

    void resume_game() {
        game_timer.resume();
        shot_clock.resume();
        notify_changes();
    }

    void end_game() {
        game_timer.reset();
        shot_clock.reset();
        is_game_over = true;
        notify_changes();
    }

    void next_quarter() {
        if (quarter < max_quarters) {
            quarter++;
            shot_clock.reset();
            shot_clock.start(24000ms);
            notify_changes();
        }
    }

    void reset_shot_clock() {
        shot_clock.reset();
        shot_clock.start(24000ms);
        notify_changes();
    }

    void add_time_to_shot_clock(milliseconds time) {
        shot_clock.add_time(time);
        notify_changes();
    }

    void add_time_to_game(milliseconds time) {
        game_timer.add_time(time);
        notify_changes();
    }

    void update_home_score(int points) {
        home_team.add_score(points);
        reset_shot_clock();
        notify_changes();
    }

    void update_away_score(int points) {
        away_team.add_score(points);
        reset_shot_clock();
        notify_changes();
    }

    void home_foul() {
        home_team.add_foul();
        notify_changes();
    }

    void away_foul() {
        away_team.add_foul();
        notify_changes();
    }

    void home_timeout() {
        if (home_team.use_timeout()) {
            pause_game();
            notify_changes();
        }
    }

    void away_timeout() {
        if (away_team.use_timeout()) {
            pause_game();
            notify_changes();
        }
    }

    map<string, string> get_game_state() const {
        map<string, string> state;
        state["home_name"] = home_team.get_name();
        state["away_name"] = away_team.get_name();
        state["home_score"] = to_string(home_team.get_score());
        state["away_score"] = to_string(away_team.get_score());
        state["home_fouls"] = to_string(home_team.get_fouls());
        state["away_fouls"] = to_string(away_team.get_fouls());
        state["home_timeouts"] = to_string(home_team.get_timeouts());
        state["away_timeouts"] = to_string(away_team.get_timeouts());
        state["game_time"] = game_timer.get_formatted_time();
        state["shot_clock"] = to_string(duration_cast<seconds>(shot_clock.get_remaining_time()).count());
        state["quarter"] = to_string(quarter);
        state["game_over"] = is_game_over ? "true" : "false";
        return state;
    }
};
