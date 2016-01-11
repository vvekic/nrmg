using System;
using System.Collections.Generic;
using System.Text;

namespace Heroes5RandomMap.DataProc
{
    static class Global
    {
        public static bool ArrayEquals<T>(T[] array1, int startIndex1, T[] array2 , int startIndex2, int Length)//测试两个数组是否匹配
        {
            if (array1 == array2) return true;
            if (array1 == null || array2 == null) return false;
            if (array1.Length != array1.Length) return false;
            if (array1.Length < (startIndex1+Length) || array2.Length < (startIndex2+Length)) return false;
            /* 以下复杂部分暂时不用
            Type equatable = typeof(T).GetInterface("IEquatable`1");
            if (equatable != null)
            {
                MethodInfo equalMethod = equatable.GetMethod("Equals", BindingFlags.Instance | BindingFlags.Public);
                for (int i = 0; i < Length; i++)
                {
                    if (!(bool)equalMethod.Invoke(array1[i+startIndex1], new object[] { array2[i+startIndex2] })) return false;
                }
                return true;
            }
            Type comparable = typeof(T).GetInterface("IComparable`1");
            if (comparable != null)
            {
                MethodInfo compareMethod = comparable.GetMethod("CompareTo", BindingFlags.Instance | BindingFlags.Public);
                for (int i = 0; i < Length; i++)
                {
                    if (!(bool)compareMethod.Invoke(array1[i + startIndex1], new object[] { array2[i + startIndex2] })) return false;
                }
                return true;
            }
            */
            for (int i = 0; i < Length; i++)
            {
                if (!array1[i+startIndex1].Equals(array2[i+startIndex2])) return false;
            }
            return true;
        }

        public static bool ArrayEquals<T>(T[] array1, int startIndex1, T[] array2) //测试Array1是否匹配整个array2
        {
            return ArrayEquals<T>(array1, startIndex1, array2, 0, array2.Length);
        }

        public static string ByteArrayToString(byte[] bArray)
        {
            string result = "";
            foreach (byte aByte in bArray)
            {
                result += aByte.ToString("X2");
            }
            return result;
        }
    }
}
