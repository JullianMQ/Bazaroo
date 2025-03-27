import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'cart.dart';
import 'productsPage.dart';
import 'address.dart';
import 'payment.dart';
import 'history.dart';
import 'nav/customer_nav.dart';  
import 'logincustomer.dart';
import 'buyerProfile.dart';
import 'productDetails.dart';

class HomeScreen extends StatefulWidget {
  final String userId;
  const HomeScreen({Key? key, required this.userId}) : super(key: key);

  @override
  _HomeScreenState createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  List products = [];
  final String baseUrl = "http://localhost:3000";

   @override
  void initState() {
    super.initState();
    fetchProducts();
  }
  @override

  Future<void> fetchProducts() async {
  final url = Uri.parse('$baseUrl/v1/products');
    try {
      final response = await http.get(url);
      if (response.statusCode == 200) {
        List jsonData = json.decode(response.body);
        setState(() {
          products = jsonData;
        });
      } else {
        print('Failed to load products: ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching products: $e');
    }
  }


  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        automaticallyImplyLeading: false,
        toolbarHeight: 80,
        actions: <Widget>[Container()],
        iconTheme: const IconThemeData(color: Colors.white),
        backgroundColor: Colors.red,
        title: Row(
          children: [
            GestureDetector(
              onTap: () {
                Navigator.pushReplacement(
                  context,
                  MaterialPageRoute(
                    builder: (context) => HomeScreen(userId: widget.userId)
                  ),
                );
              },
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 8.0),
                child: Row(
                  children: [
                    Image.asset(
                      'assets/Bazaroo-logo-white.png',
                      height: 40,
                    ),
                    Image.asset(
                      'assets/Bazaroo-text-white.png',
                      height: 30,
                    ),
                  ],
                ),
              ),
            ),
            const Spacer(),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8.0),
              child: Builder(
                builder: (context) => IconButton(
                  icon: const Icon(Icons.shopping_cart, color: Colors.white, size: 25),
                  onPressed: () => Navigator.push(
                    context,
                    MaterialPageRoute(builder: (context) => Cart(userId: widget.userId)),
                  ),
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8.0),
              child: Builder(
                builder: (context) => IconButton(
                  icon: const Icon(Icons.menu, color: Colors.white, size: 25),
                  onPressed: () => Scaffold.of(context).openEndDrawer(),
                ),
              ),
            ),
          ],
        ),
      ),
      endDrawer: Drawer(
        child: ListView(
          padding: EdgeInsets.zero,
          children: [
            SizedBox(
              height: 90,
              child: DrawerHeader(
                child: Align(
                  alignment: Alignment.topRight,
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 10.0, vertical: 5.0),
                    child: IconButton(
                      icon: const Icon(Icons.close, color: Colors.red),
                      onPressed: () => Navigator.of(context).pop(),
                    ),
                  ),
                ),
              ),
            ),
            ListTile(title: const Text('My Account'), onTap: () {
              Navigator.push(context, MaterialPageRoute(builder: (context) => BuyerProfile(userId: widget.userId)));
            }),
            
            ListTile(title: const Text('Purchase History'), onTap: () {
              Navigator.push(context, MaterialPageRoute(builder: (context) => PurchaseHistory(userId: widget.userId)));
            }),
            ListTile(title: const Text('Addresses'), onTap: () {
              Navigator.push(context, MaterialPageRoute(builder: (context) => AddressScreen(userId: widget.userId)));
            }),
            ListTile(title: const Text('Payment Options'), onTap: () {
              Navigator.push(context, MaterialPageRoute(builder: (context) => PaymentOptions(userId: widget.userId)));
            }),
            const ListTile(title: Text('Change Password')),
            const ListTile(title: Text('Settings')),
            ListTile(title: const Text('Log out'),  onTap: () {
              Navigator.pushReplacement(context, MaterialPageRoute(builder: (context) => LoginCustomer()),
              );
            }),
          ],
        ),
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(20.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Featured
              Container(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Container( // Heading
                      margin: const EdgeInsets.only(bottom: 10),
                      child: const Text('Featured',
                        style: TextStyle(
                          fontSize: 24,
                          fontWeight: FontWeight.bold,
                          color: Colors.red,
                        ),
                      ),
                    ),
                    Container( // Contents or image placeholder
                      margin: const EdgeInsets.only(bottom: 40),
                      child: AspectRatio(
                        aspectRatio: 16 / 9,
                        child: Container(
                          decoration: BoxDecoration(
                            color: Colors.grey[300],
                            borderRadius: BorderRadius.circular(15),
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              // Products
              Container(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Container( // Heading
                      margin: const EdgeInsets.only(bottom: 10),
                      child: GestureDetector(
                        onTap: () {
                          Navigator.push(
                            context,
                            MaterialPageRoute(builder: (context) => ProductsPage(userId: widget.userId)),
                          );
                        },
                        child: const Text(
                          'Products',
                          style: TextStyle(
                            fontSize: 24,
                            fontWeight: FontWeight.bold,
                            color: Colors.red,
                          ),
                        ),
                      ),
                    ),
                    Container( // Contents
                      child: GridView.builder(
                        shrinkWrap: true,
                        physics: const NeverScrollableScrollPhysics(),
                        gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                          crossAxisCount: 2,
                          crossAxisSpacing: 16,
                          mainAxisSpacing: 16,
                          childAspectRatio: 1,
                        ),
                        itemCount: products.length,
                        itemBuilder: (context, index) {
                          final product = products[index];
                          String imageUrl = "$baseUrl${product['prod_image']['String']}";
                          return GestureDetector(
                            onTap: () {
                              Navigator.push(
                                context,
                                MaterialPageRoute(
                                  builder: (context) => ProductDetailScreen(
                                    userId: widget.userId,
                                    prodId: product['prod_id']
                                  ),
                                ),
                              );
                            },
                            child: Container(
                              decoration: BoxDecoration(
                                borderRadius: BorderRadius.circular(15),
                                image: DecorationImage(
                                  image: NetworkImage(imageUrl),
                                  fit: BoxFit.cover,
                                ),
                              ),
                            ),
                          );
                        },
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
      bottomNavigationBar: BottomNavBar(userId: widget.userId),
    );
  }
}

