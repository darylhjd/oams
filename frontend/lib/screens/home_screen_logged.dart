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
      ],
    );
  }

  Widget _desktop() {
    return ListView(
      padding: const EdgeInsets.all(_desktopPadding),
      children: [
        _UpcomingSessionsCalendar(),
      ],
    );
  }
}

final _eventsMapProvider =
    Provider<Map<String, List<UpcomingClassGroupSession>>>((ref) {
  final Map<String, List<UpcomingClassGroupSession>> eventsMap = {};
  final upcomingSessions =
      ref.watch(sessionUserProvider).requireValue.upcomingSessions;

  for (var element in upcomingSessions) {
    eventsMap.update(
      _UpcomingSessionsCalendarState.dateComparator.format(element.startTime),
      (value) {
        value.add(element);
        return value;
      },
      ifAbsent: () => [element],
    );
  }

  return eventsMap;
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

  // We use the string formatting of the date as the key for the selected events
  // map since we need to compare within a specific timezone and not UTC time.
  // If not, we may get an edge case where days which are the same in a particular
  // timezone are not the same in UTC time. However, the dart implementation
  // for getting information from a DateTime object works on UTC time. For example,
  // a time of 7.30+0800 will have its day field be one day before 8.00+0800.
  static final DateFormat dateComparator = DateFormat("yyyy-MM-dd");

  late List<UpcomingClassGroupSession> _selectedEvents;
  CalendarFormat _calendarFormat = CalendarFormat.month;
  DateTime _focusedDay = DateTime.now();
  DateTime _selectedDay = DateTime.now();

  @override
  void initState() {
    super.initState();
    _selectedEvents = _getEventsForDay(DateTime.now());
  }

  @override
  Widget build(BuildContext context) {
    return TableCalendar(
      firstDay: DateTime.now(),
      lastDay: DateTime.now().add(_dateBuffer),
      focusedDay: _focusedDay,
      calendarFormat: _calendarFormat,
      weekNumbersVisible: true,
      selectedDayPredicate: (day) => _isSameDay(day, _selectedDay),
      onFormatChanged: (format) {
        setState(() {
          _calendarFormat = format;
        });
      },
      onDaySelected: (selectedDay, focusedDay) {
        if (!_isSameDay(selectedDay, _selectedDay)) {
          setState(() {
            _selectedDay = selectedDay;
            _focusedDay = focusedDay;
            _selectedEvents = _getEventsForDay(selectedDay);
          });
        }
      },
      eventLoader: (day) => _getEventsForDay(day),
    );
  }

  // Get all events happening on a particular day.
  List<UpcomingClassGroupSession> _getEventsForDay(DateTime day) {
    return ref.read(_eventsMapProvider)[dateComparator.format(day)] ?? [];
  }

  static bool _isSameDay(DateTime d1, DateTime d2) {
    return dateComparator.format(d1) == dateComparator.format(d2);
  }
}
