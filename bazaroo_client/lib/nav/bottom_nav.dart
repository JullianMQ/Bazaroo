import 'package:flutter/material.dart';
import '../main.dart';

class BottomNavBar extends StatelessWidget {

  @override
  Widget build(BuildContext context) {
    return BottomNavigationBar(
       backgroundColor: Color(0xFFE5E7EB),
      selectedItemColor: Colors.red,
      unselectedItemColor: Colors.red,
      showSelectedLabels: false,
      showUnselectedLabels: false,
      items: [
        BottomNavigationBarItem(
          icon: GestureDetector(
            onTap: () {
              Navigator.pushReplacement(context,
                MaterialPageRoute(builder: (context) => MyApp()),
              );
            },
            child: Icon(Icons.home, size: 30,),
          ),
          label: 'Home',
        ),
        BottomNavigationBarItem(
          icon: GestureDetector(
            onTap: () {
            },
            child: Icon(Icons.notifications, size: 30),
          ),
          label: 'Notifications',
        ),
        BottomNavigationBarItem(
          icon: GestureDetector(
            onTap: () {
            },
            child: Icon(Icons.mail, size: 30),
          ),
          label: 'Mail',
        ),
      ],
    );
  }
}
