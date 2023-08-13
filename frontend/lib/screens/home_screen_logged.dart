import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/api/models.dart';
import 'package:frontend/providers/session.dart';
import 'package:frontend/screens/screen_template.dart';
import 'package:intl/intl.dart';
import 'package:responsive_framework/responsive_breakpoints.dart';
import 'package:table_calendar/table_calendar.dart';

// Shows the logged-in version of the home screen.
class HomeScreenLoggedIn extends ConsumerWidget {
  static const double _mobilePadding = 10;
  static const double _desktopPadding = 20;
  static const double _desktopMaxHeight = 400;

  const HomeScreenLoggedIn({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return ScreenTemplate(
      ResponsiveBreakpoints.of(context).isMobile ? _mobile() : _desktop(),
    );
  }

  Widget _mobile() {
    return ListView(
      padding: const EdgeInsets.all(_mobilePadding),
      children: [
        _UpcomingSessionsCalendar(),
        const SizedBox(height: _mobilePadding),
        const _SelectedDaySessionsPreviewer(true),
        const SizedBox(height: _mobilePadding),
        const Placeholder(),
      ],
    );
  }

  Widget _desktop() {
    return ListView(
      padding: const EdgeInsets.all(_desktopPadding),
      children: [
        SizedBox(
          height: _desktopMaxHeight,
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Flexible(child: _UpcomingSessionsCalendar()),
              const _SelectedDaySessionsPreviewer(false),
            ],
          ),
        ),
        const SizedBox(height: _desktopPadding),
        const Placeholder(),
      ],
    );
  }
}

// The notifier for getting the current selected day events of the calendar.
class _SelectedDayEventsNotifier
    extends StateNotifier<List<UpcomingClassGroupSession>> {
  // We use the string formatting of the date as the key for the selected events
  // map since we need to compare within a specific timezone and not UTC time.
  // If not, we may get an edge case where days which are the same in a particular
  // timezone are not the same in UTC time. However, the dart implementation
  // for getting information from a DateTime object works on UTC time. For example,
  // a time of 7.30+0800 will have its day field be one day before 8.00+0800.
  static final DateFormat _dateComparator = DateFormat("yyyy-MM-dd");

  final Map<String, List<UpcomingClassGroupSession>> _eventsMap = {};

  _SelectedDayEventsNotifier(List<UpcomingClassGroupSession> upcomingSessions)
      : super([]) {
    for (var element in upcomingSessions) {
      _eventsMap.update(
        _dateComparator.format(element.startTime),
        (value) {
          value.add(element);
          return value;
        },
        ifAbsent: () => [element],
      );
    }

    state = this[DateTime.now()] ?? [];
  }

  // Override the index operator.
  operator [](DateTime d) => _eventsMap[_dateComparator.format(d)] ?? [];

  void setSelectedDayEvents(DateTime d) {
    state = this[d];
  }

  static bool isSameDay(DateTime d1, DateTime d2) =>
      _dateComparator.format(d1) == _dateComparator.format(d2);
}

// This provider provides the current selected day events.
final _selectedDayEventsProvider = StateNotifierProvider<
    _SelectedDayEventsNotifier, List<UpcomingClassGroupSession>>((ref) {
  final upcomingSessions =
      ref.watch(sessionUserProvider).requireValue.upcomingSessions;
  return _SelectedDayEventsNotifier(upcomingSessions);
});

// This provides a calendar view that shows all upcoming class group sessions
// for a user.
class _UpcomingSessionsCalendar extends ConsumerStatefulWidget {
  @override
  _UpcomingSessionsCalendarState createState() =>
      _UpcomingSessionsCalendarState();
}

// This stores the state for the upcoming class group session calendar.
class _UpcomingSessionsCalendarState extends ConsumerState {
  static const Duration _dateBuffer = Duration(days: 6 * 31);

  CalendarFormat _calendarFormat = CalendarFormat.month;
  DateTime _focusedDay = DateTime.now();
  DateTime _selectedDay = DateTime.now();

  @override
  Widget build(BuildContext context) {
    return TableCalendar(
      firstDay: DateTime.now(),
      lastDay: DateTime.now().add(_dateBuffer),
      focusedDay: _focusedDay,
      calendarFormat: _calendarFormat,
      weekNumbersVisible: true,
      formatAnimationCurve: Curves.easeInOutCubic,
      availableGestures: AvailableGestures.horizontalSwipe,
      selectedDayPredicate: (day) =>
          _SelectedDayEventsNotifier.isSameDay(day, _selectedDay),
      onFormatChanged: (format) {
        setState(() {
          _calendarFormat = format;
        });
      },
      onDaySelected: (selectedDay, focusedDay) {
        if (!_SelectedDayEventsNotifier.isSameDay(selectedDay, _selectedDay)) {
          setState(() {
            _selectedDay = selectedDay;
            _focusedDay = focusedDay;
            ref
                .watch(_selectedDayEventsProvider.notifier)
                .setSelectedDayEvents(selectedDay);
          });
        }
      },
      eventLoader: (day) => ref.watch(_selectedDayEventsProvider.notifier)[day],
    );
  }
}

// This provides the mapping from a class type to the color coding.
const Map<ClassType, Color> _colorMap = {
  ClassType.lec: Colors.lightBlueAccent,
  ClassType.tut: Colors.lightGreen,
  ClassType.lab: Colors.orangeAccent,
};

// Shows the selected day's sessions.
class _SelectedDaySessionsPreviewer extends ConsumerWidget {
  static const Text _header = Text(
    "Selected date sessions",
    textAlign: TextAlign.center,
  );
  static const Text _noEvents = Text(
    "No classes on this date. Hooray!",
    textAlign: TextAlign.center,
  );

  static const double _mobileMaxHeight = 400;
  static const double _mobilePadding = 5;

  static const double _desktopPadding = 10;
  static const double _desktopWidth = 400;

  final bool _isMobile;

  const _SelectedDaySessionsPreviewer(this._isMobile);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return _isMobile ? _mobile(context, ref) : _desktop(context, ref);
  }

  Widget _mobile(BuildContext context, WidgetRef ref) {
    final previews = _getEventPreviews(ref);

    return ConstrainedBox(
      constraints: const BoxConstraints(maxHeight: _mobileMaxHeight),
      child: Card(
        child: Container(
          padding: const EdgeInsets.all(_mobilePadding),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            mainAxisSize: MainAxisSize.min,
            children: [
              _header,
              const Divider(),
              previews.isEmpty
                  ? _noEvents
                  : ListView(
                      shrinkWrap: true,
                      children: previews,
                    ),
              const Divider(),
              const _SelectedDaySessionsPreviewerFooter(),
            ],
          ),
        ),
      ),
    );
  }

  Widget _desktop(BuildContext context, WidgetRef ref) {
    final previews = _getEventPreviews(ref);

    return SizedBox(
      width: _desktopWidth,
      child: Card(
        child: Container(
          padding: const EdgeInsets.all(_desktopPadding),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              _header,
              const Divider(),
              Expanded(
                child: previews.isEmpty
                    ? Container(
                        alignment: Alignment.center,
                        child: _noEvents,
                      )
                    : ListView(
                        children: previews,
                      ),
              ),
              const Divider(),
              const _SelectedDaySessionsPreviewerFooter(),
            ],
          ),
        ),
      ),
    );
  }

  List<_EventPreview> _getEventPreviews(WidgetRef ref) {
    return ref
        .watch(_selectedDayEventsProvider)
        .map((e) => _EventPreview(e, _isMobile))
        .toList();
  }
}

// Provides a color legend on the footer of the sessions previewer.
class _SelectedDaySessionsPreviewerFooter extends StatelessWidget {
  const _SelectedDaySessionsPreviewerFooter();

  @override
  Widget build(BuildContext context) {
    final children = _colorMap.entries
        .map(
          (e) => Padding(
            padding: const EdgeInsets.symmetric(horizontal: 10),
            child: Text.rich(
              TextSpan(
                children: [
                  WidgetSpan(
                    alignment: PlaceholderAlignment.middle,
                    child: Icon(
                      Icons.circle,
                      color: e.value,
                      size: 10,
                    ),
                  ),
                  TextSpan(text: e.key.name),
                ],
              ),
            ),
          ),
        )
        .toList();

    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: children,
    );
  }
}

// Shows information for one upcoming session.
class _EventPreview extends StatelessWidget {
  static const double _mobilePadding = 10;
  static const double _desktopPadding = 20;

  final UpcomingClassGroupSession _session;
  final bool _isMobile;

  const _EventPreview(this._session, this._isMobile);

  @override
  Widget build(BuildContext context) {
    final timeFormatter = DateFormat("H:mm");

    return Card(
      color: _colorMap[_session.classType],
      child: Container(
        padding: EdgeInsets.all(_isMobile ? _mobilePadding : _desktopPadding),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              "${_session.code} ${_session.name}",
              style: Theme.of(context)
                  .textTheme
                  .bodyMedium
                  ?.copyWith(fontWeight: FontWeight.bold),
            ),
            Text(
              _session.classType.name,
              style: Theme.of(context).textTheme.bodyMedium,
            ),
            Text(
              "${timeFormatter.format(_session.startTime)} - ${timeFormatter.format(_session.endTime)}",
              style: Theme.of(context)
                  .textTheme
                  .bodyMedium
                  ?.copyWith(fontStyle: FontStyle.italic),
            ),
          ],
        ),
      ),
    );
  }
}
