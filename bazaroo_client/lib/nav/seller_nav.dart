import 'package:flutter/material.dart';
import '../homeseller.dart';

class BottomNavBar extends StatelessWidget {
  final String userId;
  const BottomNavBar({Key? key, required this.userId}) : super(key: key);

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
                MaterialPageRoute(builder: (context) => SellerHome(userId: userId)),
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
